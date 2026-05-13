package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/errs"
	"example.com/taskservice/internal/types"

	"github.com/go-playground/validator/v10"
)

type Scheduler interface {
	IsDue(rt taskdomain.RecurringTask, now time.Time) (time.Time, bool, error)
	ValidateConfig(raw json.RawMessage) error
}

type IntervalScheduler struct{ Validate *validator.Validate }

func (s *IntervalScheduler) calculateNext(rt taskdomain.RecurringTask, now time.Time) (time.Time, error) {
	cfg, err := s.parseConfig(rt)
	if err != nil {
		return time.Time{}, err
	}

	var last time.Time
	if rt.LastRunAt == nil {
		last = rt.StartDate
	} else {
		last = *rt.LastRunAt
	}

	if !last.Before(now) {
		return last, nil
	}

	for {
		next, errLoop := calculateNextForIntervalConfig(cfg, last)
		if errLoop != nil {
			return time.Time{}, errLoop
		}
		if !next.Before(now) {
			return next, nil
		}
		last = next
	}
}

func (s *IntervalScheduler) IsDue(rt taskdomain.RecurringTask, now time.Time) (time.Time, bool, error) {
	// Check if LastRunAt is in the past
	if rt.LastRunAt != nil && !rt.LastRunAt.Before(now) {
		return time.Time{}, false, nil
	}

	next, err := s.calculateNext(rt, now)
	if err != nil {
		return time.Time{}, false, err
	}
	// Check if next is after the EndDate
	if rt.EndDate != nil && next.After(*rt.EndDate) {
		return time.Time{}, false, nil
	}

	return next, true, nil
}

func (s *IntervalScheduler) parseConfig(rt taskdomain.RecurringTask) (IntervalConfig, error) {
	var cfg IntervalConfig
	decoder := json.NewDecoder(bytes.NewBuffer(rt.Config))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&cfg); err != nil {
		return IntervalConfig{}, errors.Join(errs.ErrInvalidInput, err)
	}
	return cfg, nil
}

func (s *IntervalScheduler) ValidateConfig(raw json.RawMessage) error {
	cfg, err := parseRawConfig[IntervalConfig](raw)
	if err != nil {
		return err
	}

	if err := s.Validate.Struct(cfg); err != nil {
		return errors.Join(errs.ErrInvalidInput, err)
	}

	switch cfg.Frequency {
	case types.FrequencyDaily:
		if cfg.Interval > 365 {
			return errors.New("daily interval too large")
		}

	case types.FrequencyWeekly:
		if cfg.Interval > 52 {
			return errors.New("weekly interval too large")
		}

	case types.FrequencyMonthly:
		if cfg.Interval > 12 {
			return errors.New("monthly interval too large")
		}

	case types.FrequencyYearly:
		if cfg.Interval > 10 {
			return errors.New("yearly interval too large")
		}
	}

	return nil
}

type EvenOddDaysScheduler struct{ Validate *validator.Validate }

func (s *EvenOddDaysScheduler) calculateNext(rt taskdomain.RecurringTask, now time.Time) (time.Time, error) {
	cfg, err := s.parseConfig(rt)
	if err != nil {
		return time.Time{}, err
	}

	var current time.Time

	if rt.LastRunAt == nil {
		current = rt.StartDate
	} else {
		current = rt.LastRunAt.AddDate(0, 0, 1)
	}

	for {
		day := current.Day()

		switch cfg.Mode {
		case "odd":
			if day%2 == 1 {
				if !current.Before(now) {
					return current, nil
				}
			}

		case "even":
			if day%2 == 0 {
				if !current.Before(now) {
					return current, nil
				}
			}

		default:
			return time.Time{}, errors.New("unknown mode")
		}

		current = current.AddDate(0, 0, 1)
	}
}

func (s *EvenOddDaysScheduler) parseConfig(rt taskdomain.RecurringTask) (OddEvenConfig, error) {
	var cfg OddEvenConfig
	decoder := json.NewDecoder(bytes.NewBuffer(rt.Config))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&cfg); err != nil {
		return OddEvenConfig{}, errors.Join(errs.ErrInvalidInput, err)
	}
	return cfg, nil
}

func (s *EvenOddDaysScheduler) IsDue(rt taskdomain.RecurringTask, now time.Time) (time.Time, bool, error) {
	if rt.LastRunAt != nil && !rt.LastRunAt.Before(now) {
		return time.Time{}, false, nil
	}

	next, err := s.calculateNext(rt, now)
	if err != nil {
		return time.Time{}, false, err
	}

	if rt.EndDate != nil && next.After(*rt.EndDate) {
		return time.Time{}, false, nil
	}

	return next, true, nil
}

func (s *EvenOddDaysScheduler) ValidateConfig(raw json.RawMessage) error {
	cfg, err := parseRawConfig[OddEvenConfig](raw)
	if err != nil {
		return err
	}

	if err := s.Validate.Struct(cfg); err != nil {
		return errors.Join(errs.ErrInvalidInput, err)
	}
	return nil
}

type SpecificDateScheduler struct{ Validate *validator.Validate }

func (s *SpecificDateScheduler) calculateNext(rt taskdomain.RecurringTask, now time.Time) (time.Time, error) {
	cfg, err := s.parseConfig(rt)
	if err != nil {
		return time.Time{}, err
	}

	var last time.Time
	if rt.LastRunAt == nil {
		last = rt.StartDate
	} else {
		last = *rt.LastRunAt
	}

	var next time.Time

	for _, date := range cfg.Dates {
		// skip past dates
		if date.Before(now) {
			continue
		}

		// skip already consumed dates
		if date.Before(last) {
			continue
		}

		// find minimum suitable date
		if next.IsZero() || date.Before(next) {
			next = date
		}
	}

	if next.IsZero() {
		return time.Time{}, errors.New("no suitable next date")
	}

	return next, nil
}

func (s *SpecificDateScheduler) parseConfig(rt taskdomain.RecurringTask) (SpecificDateConfig, error) {
	var cfg SpecificDateConfig
	decoder := json.NewDecoder(bytes.NewBuffer(rt.Config))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&cfg); err != nil {
		return SpecificDateConfig{}, errors.Join(errs.ErrInvalidInput, err)
	}
	return cfg, nil
}

func (s *SpecificDateScheduler) IsDue(rt taskdomain.RecurringTask, now time.Time) (time.Time, bool, error) {
	if rt.LastRunAt != nil && !rt.LastRunAt.Before(now) {
		return time.Time{}, false, nil
	}

	next, err := s.calculateNext(rt, now)
	if err != nil {
		return time.Time{}, false, err
	}

	if rt.EndDate != nil && next.After(*rt.EndDate) {
		return time.Time{}, false, nil
	}

	return next, true, nil
}

func (s *SpecificDateScheduler) ValidateConfig(raw json.RawMessage) error {
	cfg, err := parseRawConfig[SpecificDateConfig](raw)
	if err != nil {
		return err
	}

	if err := s.Validate.Struct(cfg); err != nil {
		return errors.Join(errs.ErrInvalidInput, err)
	}
	return nil
}
