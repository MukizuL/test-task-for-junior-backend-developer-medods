package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/errs"
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
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("unexpected trailing data")
	}
	return nil
}

func writeUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, taskdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, errs.ErrInvalidInput):
		uw, ok := err.(interface{ Unwrap() []error })
		if !ok {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if validationErrs := getValidationErrors(uw.Unwrap()[1]); len(validationErrs) > 0 {
			writeValidationError(w, http.StatusBadRequest, validationErrs)
			return
		}
		if syntaxErr, ok := errors.AsType[*json.SyntaxError](uw.Unwrap()[1]); ok {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid json syntax at position %d", syntaxErr.Offset))
			return
		}
		if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](uw.Unwrap()[1]); ok {
			writeError(w, http.StatusBadRequest, fmt.Errorf("field %q must be %s", typeErr.Field, typeErr.Type.String()))
			return
		}
		writeError(w, http.StatusBadRequest, uw.Unwrap()[1])
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
