package task

import (
	"context"
	"fmt"

	"example.com/taskservice/internal/domain/taskdomain"
)

func (s *Service) CreateRecurringTask(ctx context.Context, input CreateRecurringInput) (*taskdomain.RecurringTask, error) {
	normalized, err := validateCreateRecurringInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.RecurringTask{
		Title:       normalized.Title,
		Description: normalized.Description,
		Frequency:   normalized.Frequency,
		Interval:    normalized.Interval,
		StartDate:   normalized.StartDate,
		EndDate:     normalized.EndDate,
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

	normalized, err := validateUpdateRecurringInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.RecurringTask{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Frequency:   normalized.Frequency,
		Interval:    normalized.Interval,
		StartDate:   normalized.StartDate,
		EndDate:     normalized.EndDate,
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
