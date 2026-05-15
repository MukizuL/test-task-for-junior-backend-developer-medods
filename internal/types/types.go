package types

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
	FrequencyYearly  Frequency = "yearly"
)

type RecurrenceType string

const (
	RecurrenceInterval      RecurrenceType = "interval"
	RecurrenceEvenOddDays   RecurrenceType = "even_odd_days"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)
