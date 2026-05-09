package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/errs"
	"example.com/taskservice/internal/types"

	"github.com/go-playground/validator/v10"
)

type Scheduler interface {
	IsDue(rt taskdomain.RecurringTask, now time.Time) (bool, error)
	NextRun(rt taskdomain.RecurringTask) (time.Time, error)
	ValidateConfig(raw json.RawMessage) error
}

type IntervalScheduler struct{ Validate *validator.Validate }

func (s *IntervalScheduler) calculateNext(rt taskdomain.RecurringTask) (time.Time, error) {
	last := *rt.LastRunAt

	cfg, err := s.parseConfig(rt)
	if err != nil {
		return time.Time{}, err
	}

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

func (s *IntervalScheduler) parseConfig(rt taskdomain.RecurringTask) (IntervalConfig, error) {
	var cfg IntervalConfig
	decoder := json.NewDecoder(bytes.NewBuffer(rt.Config))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&cfg); err != nil {
		return IntervalConfig{}, errors.Join(errs.ErrInvalidInput, err)
	}
	return cfg, nil
}

func (s *IntervalScheduler) IsDue(rt taskdomain.RecurringTask, now time.Time) (bool, error) {
	if rt.LastRunAt == nil {
		return !rt.StartDate.After(now), nil
	}

	next, err := s.calculateNext(rt)
	if err != nil {
		return false, err
	}
	return !next.After(now), nil
}

func (s *IntervalScheduler) NextRun(rt taskdomain.RecurringTask) (time.Time, error) {
	if rt.LastRunAt == nil {
		return rt.StartDate, nil
	}
	return s.calculateNext(rt)
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

type OddDaysScheduler struct{ Validate *validator.Validate }

func (s *OddDaysScheduler) calculateNext(rt taskdomain.RecurringTask) (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OddDaysScheduler) IsDue(rt taskdomain.RecurringTask, now time.Time) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OddDaysScheduler) NextRun(rt taskdomain.RecurringTask) (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (s *OddDaysScheduler) ValidateConfig(raw json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}

type EvenDaysScheduler struct{ Validate *validator.Validate }

func (s *EvenDaysScheduler) calculateNext(rt taskdomain.RecurringTask) (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (s *EvenDaysScheduler) IsDue(rt taskdomain.RecurringTask, now time.Time) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *EvenDaysScheduler) NextRun(rt taskdomain.RecurringTask) (time.Time, error) {
	//TODO implement me
	panic("implement me")
}
func (s *EvenDaysScheduler) ValidateConfig(raw json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}

type YearlyDateScheduler struct{ Validate *validator.Validate }

func (s *YearlyDateScheduler) calculateNext(rt taskdomain.RecurringTask) (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (s *YearlyDateScheduler) IsDue(rt taskdomain.RecurringTask, now time.Time) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *YearlyDateScheduler) NextRun(rt taskdomain.RecurringTask) (time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (s *YearlyDateScheduler) ValidateConfig(raw json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}
