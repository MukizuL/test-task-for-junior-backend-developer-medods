# OpenCode

Chat 'ChatTitle(text="Generate openapi fo…" with OpenCode, isCustom=true)' (6b1b4ebc-e708-4a59-89f4-13651c1797b6)
Context:

Messages: 2
=================================================================================================================

==== RestoredMessage #1 ====
User: 
Generate openapi for recurring tasks
Attachments:
Attachment Name: dto.go
Attachments Kind: Visible
Attachments Source: Other / RestoredChatAttachment
Attachments Text:
```go
package handlers

import (
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	DueDate     time.Time         `json:"due_date"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	DueDate     time.Time         `json:"due_date"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

type recurringTaskMutationDTO struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Frequency   taskdomain.Frequency `json:"frequency"`
	Interval    int                  `json:"interval"`
	StartDate   time.Time            `json:"start_date"`
	EndDate     *time.Time           `json:"end_date"`
}

type recurringTaskDTO struct {
	ID          int64                `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Frequency   taskdomain.Frequency `json:"frequency"`
	Interval    int                  `json:"interval"`
	StartDate   time.Time            `json:"start_date"`
	EndDate     *time.Time           `json:"end_date"`
	LastRunAt   *time.Time           `json:"last_run_at"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func newRecurringTaskDTO(task *taskdomain.RecurringTask) recurringTaskDTO {
	return recurringTaskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Frequency:   task.Frequency,
		Interval:    task.Interval,
		StartDate:   task.StartDate,
		EndDate:     task.EndDate,
		LastRunAt:   task.LastRunAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

```
Attachment Name: openapi.json
Attachments Kind: Visible
Attachments Source: Other / RestoredChatAttachment
Attachments Text:
```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "Task Service API",
    "version": "1.0.0",
    "description": "CRUD API for tasks built with clean architecture."
  },
  "servers": [
    {
      "url": "http://localhost:8080"
    }
  ],
  "tags": [
    {
      "name": "Tasks",
      "description": "Operations with tasks"
    }
  ],
  "paths": {
    "/api/v1/tasks": {
      "get": {
        "tags": ["Tasks"],
        "summary": "List tasks",
        "responses": {
          "200": {
            "description": "Task list",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": { "$ref": "#/components/schemas/Task" }
                }
              }
            }
          }
        }
      },
      "post": {
        "tags": ["Tasks"],
        "summary": "Create task",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateTaskRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Created task",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Task" }
              }
            }
          },
          "400": { "$ref": "#/components/responses/BadRequest" }
        }
      }
    },
    "/api/v1/tasks/{id}": {
      "get": {
        "tags": ["Tasks"],
        "summary": "Get task by ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "integer", "format": "int64", "minimum": 1 }
          }
        ],
        "responses": {
          "200": {
            "description": "Task",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Task" }
              }
            }
          },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      },
      "put": {
        "tags": ["Tasks"],
        "summary": "Update task",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "integer", "format": "int64", "minimum": 1 }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateTaskRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Updated task",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Task" }
              }
            }
          },
          "400": { "$ref": "#/components/responses/BadRequest" },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      },
      "delete": {
        "tags": ["Tasks"],
        "summary": "Delete task",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": { "type": "integer", "format": "int64", "minimum": 1 }
          }
        ],
        "responses": {
          "204": { "description": "Deleted" },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      }
    }
  },
  "components": {
    "responses": {
      "BadRequest": {
        "description": "Bad request",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/Error" }
          }
        }
      },
      "NotFound": {
        "description": "Resource not found",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/Error" }
          }
        }
      }
    },
    "schemas": {
      "CreateTaskRequest": {
        "type": "object",
        "required": ["title"],
        "properties": {
          "title": {
            "type": "string",
            "example": "Prepare release"
          },
          "description": {
            "type": "string",
            "example": "Collect release notes and check migrations"
          },
          "status": {
            "type": "string",
            "enum": ["new", "in_progress", "done"],
            "example": "new"
          },
          "due_date": {
            "type": "string",
            "example": "2026-04-24T00:00:00+07:00"
          }
        }
      },
      "UpdateTaskRequest": {
        "type": "object",
        "required": ["title", "status"],
        "properties": {
          "title": {
            "type": "string",
            "example": "Prepare release"
          },
          "description": {
            "type": "string",
            "example": "Collect release notes and check migrations"
          },
          "status": {
            "type": "string",
            "enum": ["new", "in_progress", "done"],
            "example": "in_progress"
          },
          "due_date": {
            "type": "string",
            "example": "2026-04-24T00:00:00+07:00"
          }
        }
      },
      "Task": {
        "type": "object",
        "properties": {
          "id": { "type": "integer", "format": "int64", "example": 1 },
          "recurring_task_id": {"type": "integer", "format": "int64", "example": 1},
          "title": { "type": "string", "example": "Prepare release" },
          "description": {
            "type": "string",
            "example": "Collect release notes and check migrations"
          },
          "status": {
            "type": "string",
            "enum": ["new", "in_progress", "done"],
            "example": "new"
          },
          "created_at": {
            "type": "string",
            "format": "date-time",
            "example": "2026-03-24T00:00:00Z"
          },
          "updated_at": {
            "type": "string",
            "format": "date-time",
            "example": "2026-03-24T00:00:00Z"
          }
        }
      },
      "Error": {
        "type": "object",
        "properties": {
          "error": { "type": "string", "example": "task not found" }
        }
      }
    }
  }
}

```
Attachment Name: recurring_task_handler.go
Attachments Kind: Visible
Attachments Source: Other / RestoredChatAttachment
Attachments Text:
```go
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
```

# ChatGPT

ChatGPT подсказал архитектуру. С ним создавал тесты для worker-а. Также проводил code review.
UPD 1: Решили сделать polymorphic scheduler для поддержки большего числа типов итеративности.
