package taskdomain

import "time"

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
	FrequencyYearly  Frequency = "yearly"
)

type RecurringTask struct {
	ID          int64
	Title       string
	Description string
	Frequency   Frequency
	Interval    int
	StartDate   time.Time
	EndDate     *time.Time
	LastRunAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (f Frequency) Valid() bool {
	switch f {
	case FrequencyDaily, FrequencyWeekly, FrequencyMonthly, FrequencyYearly:
		return true
	default:
		return false
	}
}
