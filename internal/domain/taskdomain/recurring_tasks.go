package taskdomain

import (
	"encoding/json"
	"time"

	"example.com/taskservice/internal/types"
)

type RecurringTask struct {
	ID          int64
	Title       string
	Description string
	Type        types.RecurrenceType
	Config      json.RawMessage
	StartDate   time.Time
	EndDate     *time.Time
	LastRunAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
