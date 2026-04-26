package postgres

import taskdomain "example.com/taskservice/internal/domain/taskdomain"

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

	task.Status = taskdomain.Status(status)

	return &task, nil
}

func scanRecurringTask(scanner taskScanner) (*taskdomain.RecurringTask, error) {
	var (
		task      taskdomain.RecurringTask
		frequency string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&frequency,
		&task.Interval,
		&task.StartDate,
		&task.EndDate,
		&task.LastRunAt,
		&task.CreatedAt,
	); err != nil {
		return nil, err
	}

	task.Frequency = taskdomain.Frequency(frequency)

	return &task, nil
}
