package worker

import (
	"time"

	"example.com/taskservice/internal/types"
)

type IntervalConfig struct {
	Frequency types.Frequency `json:"frequency" validate:"required,oneof=daily weekly monthly yearly"`
	Interval  int             `json:"interval" validate:"required,min=1"`
}

type OddEvenConfig struct {
	Mode string `json:"mode" validate:"required,oneof=odd even"`
}

type SpecificDateConfig struct {
	Dates []time.Time `json:"dates" validate:"required"`
}
