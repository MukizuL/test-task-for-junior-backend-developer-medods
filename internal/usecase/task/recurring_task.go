package task

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"example.com/taskservice/internal/domain/taskdomain"
)

func (s *Service) CreateRecurringTask(ctx context.Context, input CreateRecurringInput) (*taskdomain.RecurringTask, error) {
	err := s.validate.Struct(input)
	if err != nil {
		return nil, errors.Join(ErrInvalidInput, err)
	}

	model := &taskdomain.RecurringTask{
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Frequency:   input.Frequency,
		Interval:    input.Interval,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}

	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.CreateRecurringTask(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetRecurringTaskByID(ctx context.Context, id int64) (*taskdomain.RecurringTask, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetRecurringTaskByID(ctx, id)
}

func (s *Service) UpdateRecurringTask(ctx context.Context, id int64, input UpdateRecurringInput) (*taskdomain.RecurringTask, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	err := s.validate.Struct(input)
	if err != nil {
		return nil, errors.Join(ErrInvalidInput, err)
	}

	model := &taskdomain.RecurringTask{
		ID:          id,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Frequency:   input.Frequency,
		Interval:    input.Interval,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		UpdatedAt:   s.now(),
	}

	updated, err := s.repo.UpdateRecurringTask(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteRecurringTask(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.DeleteRecurringTask(ctx, id)
}

func (s *Service) ListRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error) {
	return s.repo.ListRecurringTasks(ctx)
}
