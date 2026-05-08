package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"example.com/taskservice/internal/domain/taskdomain"
	taskusecase "example.com/taskservice/internal/usecase/task"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func getIDFromRequest(r *http.Request) (int64, error) {
	rawID := mux.Vars(r)["id"]
	if rawID == "" {
		return 0, errors.New("missing task id")
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid task id")
	}

	if id <= 0 {
		return 0, errors.New("invalid task id")
	}

	return id, nil
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}

func writeUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, taskdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, taskusecase.ErrInvalidInput):
		if validationErrs := getValidationErrors(err); len(validationErrs) > 0 {
			writeValidationError(w, http.StatusBadRequest, validationErrs)
			return
		}
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeValidationError(w http.ResponseWriter, status int, details map[string]string) {
	writeJSON(w, status, ValidationErrorResponse{
		Error:   "Invalid task input",
		Details: details,
	})
}

func getValidationErrors(err error) map[string]string {
	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		errorsMap := make(map[string]string)
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()
			switch tag {
			case "required":
				errorsMap[field] = "This field is required"
			case "min":
				errorsMap[field] = fmt.Sprintf("Minimum length is %s", fe.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("Maximum length is %s", fe.Param())
			case "oneof":
				errorsMap[field] = fmt.Sprintf("Must be one of: %s", fe.Param())
			case "gt":
				errorsMap[field] = fmt.Sprint("Due date must be in the future")
			case "gtfield":
				errorsMap[field] = fmt.Sprint("End date must be after start date")
			default:
				errorsMap[field] = fmt.Sprintf("Failed on '%s' validation", tag)
			}
		}
		return errorsMap
	}
	return nil
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
