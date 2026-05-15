package task

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/errs"
)

func (s *Service) CreateRecurringTask(ctx context.Context, input CreateRecurringInput) (*taskdomain.RecurringTask, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	err := s.validate.Struct(input)
	if err != nil {
		return nil, errors.Join(errs.ErrInvalidInput, err)
	}

	scheduler, ok := s.schedulers[input.Type]
	if !ok {
		return nil, errors.New("unknown schedule type")
	}

	err = scheduler.ValidateConfig(input.Config)
	if err != nil {
		return nil, errors.Join(errs.ErrInvalidInput, err)
	}

	model := &taskdomain.RecurringTask{
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Config:      input.Config,
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
		return nil, fmt.Errorf("%w: id must be positive", errs.ErrInvalidInput)
	}

	return s.repo.GetRecurringTaskByID(ctx, id)
}

func (s *Service) UpdateRecurringTask(ctx context.Context, id int64, input UpdateRecurringInput) (*taskdomain.RecurringTask, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", errs.ErrInvalidInput)
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	err := s.validate.Struct(input)
	if err != nil {
		return nil, errors.Join(errs.ErrInvalidInput, err)
	}

	scheduler, ok := s.schedulers[input.Type]
	if !ok {
		return nil, errors.New("unknown schedule type")
	}

	err = scheduler.ValidateConfig(input.Config)
	if err != nil {
		return nil, errors.Join(errs.ErrInvalidInput, err)
	}

	model := &taskdomain.RecurringTask{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Config:      input.Config,
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
		return fmt.Errorf("%w: id must be positive", errs.ErrInvalidInput)
	}

	return s.repo.DeleteRecurringTask(ctx, id)
}

func (s *Service) ListRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error) {
	return s.repo.ListRecurringTasks(ctx)
}
