package worker

import "example.com/taskservice/internal/types"

type IntervalConfig struct {
	Frequency types.Frequency `json:"frequency" validate:"required,oneof=daily weekly monthly yearly"`
	Interval  int             `json:"interval" validate:"required,min=1"`
}

type OddEvenConfig struct {
	Mode string `json:"mode"` // odd/even
}

type YearlyDateConfig struct {
	Month int `json:"month"`
	Day   int `json:"day"`
}
