package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"example.com/taskservice/internal/errs"
	"example.com/taskservice/internal/types"
)

func parseRawConfig[T IntervalConfig | OddEvenConfig | YearlyDateConfig](raw json.RawMessage) (T, error) {
	var cfg T
	decoder := json.NewDecoder(bytes.NewBuffer(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&cfg); err != nil {
		return cfg, errors.Join(errs.ErrInvalidInput, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return cfg, errors.New("unexpected trailing data")
	}
	return cfg, nil
}

func calculateNextForIntervalConfig(cfg IntervalConfig, last time.Time) (time.Time, error) {
	switch cfg.Frequency {
	case types.FrequencyDaily:
		return last.AddDate(0, 0, cfg.Interval), nil
	case types.FrequencyWeekly:
		return last.AddDate(0, 0, 7*cfg.Interval), nil
	case types.FrequencyMonthly:
		return last.AddDate(0, cfg.Interval, 0), nil
	case types.FrequencyYearly:
		return last.AddDate(cfg.Interval, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("%w: %s", errUnknownFrequency, cfg.Frequency)
	}
}
