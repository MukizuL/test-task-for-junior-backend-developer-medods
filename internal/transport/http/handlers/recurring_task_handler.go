package handlers

import (
	"net/http"

	taskusecase "example.com/taskservice/internal/usecase/task"
)

func (h *TaskHandler) CreateRecurringTask(w http.ResponseWriter, r *http.Request) {
	var req recurringTaskMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := h.usecase.CreateRecurringTask(r.Context(), taskusecase.CreateRecurringInput{
		Title:       req.Title,
		Description: req.Description,
		Frequency:   req.Frequency,
		Interval:    req.Interval,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newRecurringTaskDTO(created))
}

func (h *TaskHandler) GetRecurringTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	task, err := h.usecase.GetRecurringTaskByID(r.Context(), id)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newRecurringTaskDTO(task))
}

func (h *TaskHandler) UpdateRecurringTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req recurringTaskMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	updated, err := h.usecase.UpdateRecurringTask(r.Context(), id, taskusecase.UpdateRecurringInput{
		Title:       req.Title,
		Description: req.Description,
		Frequency:   req.Frequency,
		Interval:    req.Interval,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newRecurringTaskDTO(updated))
}

func (h *TaskHandler) DeleteRecurringTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.usecase.DeleteRecurringTask(r.Context(), id); err != nil {
		writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) ListRecurringTask(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.usecase.ListRecurringTasks(r.Context())
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	response := make([]recurringTaskDTO, 0, len(tasks))
	for i := range tasks {
		response = append(response, newRecurringTaskDTO(&tasks[i]))
	}

	writeJSON(w, http.StatusOK, response)
}
