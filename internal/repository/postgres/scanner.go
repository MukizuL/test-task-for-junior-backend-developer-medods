package postgres

import (
	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/types"
)

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.RecurringTaskID,
		&task.Title,
		&task.Description,
		&status,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = types.Status(status)

	return &task, nil
}

func scanRecurringTask(scanner taskScanner) (*taskdomain.RecurringTask, error) {
	var (
		task    taskdomain.RecurringTask
		recType string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&recType,
		&task.Config,
		&task.StartDate,
		&task.EndDate,
		&task.LastRunAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Type = types.RecurrenceType(recType)

	return &task, nil
}
