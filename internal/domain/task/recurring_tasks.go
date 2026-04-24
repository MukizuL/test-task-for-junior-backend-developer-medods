package task

import "time"

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
	FrequencyYearly  Frequency = "yearly"
)

type RecurringTask struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Frequency   Frequency  `json:"frequency"`
	Interval    int        `json:"interval"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	LastRunAt   *time.Time `json:"last_run_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (f Frequency) Valid() bool {
	switch f {
	case FrequencyDaily, FrequencyWeekly, FrequencyMonthly, FrequencyYearly:
		return true
	default:
		return false
	}
}
