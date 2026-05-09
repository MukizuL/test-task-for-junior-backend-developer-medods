package worker

import "example.com/taskservice/internal/types"

type IntervalConfig struct {
	Frequency types.Frequency `json:"frequency" validate:"required,oneof=daily weekly monthly yearly"`
	Interval  int             `json:"interval" validate:"required,min=1"`
}

type OddEvenConfig struct {
	Mode string `json:"mode" validate:"required,oneof=odd even"`
}

type YearlyDateConfig struct {
	Month int `json:"month" validate:"required,min=1,max=12"`
	Day   int `json:"day" validate:"required,min=1,max=31"` // maybe implement custom rule? max should depend on days in a month
}
