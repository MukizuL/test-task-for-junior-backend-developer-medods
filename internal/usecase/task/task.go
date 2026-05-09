package task

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/errs"
)

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	err := s.validate.Struct(input)
	if err != nil {
		return nil, errors.Join(errs.ErrInvalidInput, err)
	}

	model := &taskdomain.Task{
		RecurringTaskID: nil,
		Title:           strings.TrimSpace(input.Title),
		Description:     strings.TrimSpace(input.Description),
		Status:          input.Status,
		DueDate:         input.DueDate,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", errs.ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", errs.ErrInvalidInput)
	}

	err := s.validate.Struct(input)
	if err != nil {
		return nil, errors.Join(errs.ErrInvalidInput, err)
	}

	model := &taskdomain.Task{
		ID:              id,
		RecurringTaskID: &input.RecurringTaskID,
		Title:           strings.TrimSpace(input.Title),
		Description:     strings.TrimSpace(input.Description),
		Status:          input.Status,
		DueDate:         input.DueDate,
		UpdatedAt:       s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", errs.ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}
