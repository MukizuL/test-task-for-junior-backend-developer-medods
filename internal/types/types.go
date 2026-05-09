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
	RecurrenceInterval RecurrenceType = "interval"
	RecurrenceOddDays  RecurrenceType = "odd_days"
	RecurrenceEvenDays RecurrenceType = "even_days"
	RecurrenceYearlyOn RecurrenceType = "yearly_on"
)
