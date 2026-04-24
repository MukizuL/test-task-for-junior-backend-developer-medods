package task

import (
	"fmt"
	"strings"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func validateCreateRecurringInput(input CreateRecurringInput) (CreateRecurringInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateRecurringInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Frequency.Valid() {
		return CreateRecurringInput{}, fmt.Errorf("%w: invalid frequency", ErrInvalidInput)
	}

	if input.Interval <= 0 {
		return CreateRecurringInput{}, fmt.Errorf("%w: invalid interval", ErrInvalidInput)
	}

	if input.EndDate != nil && input.StartDate.After(*input.EndDate) {
		return CreateRecurringInput{}, fmt.Errorf("%w: invalid end date", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateRecurringInput(input UpdateRecurringInput) (UpdateRecurringInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateRecurringInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Frequency.Valid() {
		return UpdateRecurringInput{}, fmt.Errorf("%w: invalid frequency", ErrInvalidInput)
	}

	if input.Interval <= 0 {
		return UpdateRecurringInput{}, fmt.Errorf("%w: invalid interval", ErrInvalidInput)
	}

	if input.EndDate != nil && input.StartDate.After(*input.EndDate) {
		return UpdateRecurringInput{}, fmt.Errorf("%w: invalid end date", ErrInvalidInput)
	}

	return input, nil
}
