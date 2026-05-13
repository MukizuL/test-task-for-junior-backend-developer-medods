package worker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"example.com/taskservice/internal/clock"
	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/types"
	"example.com/taskservice/internal/worker/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorker_IntervalScheduler(t *testing.T) {
	type want struct {
		firstRun  error
		secondRun error
	}

	tests := []struct {
		name        string
		c           *clock.FakeClock
		timeSkip    time.Duration
		mockStorage func(m *mocksWorker.MockRepo)
		want        want
	}{
		{
			name:     "Generate first and second daily task when recurring task created at now",
			c:        clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
			timeSkip: 1 * time.Hour,
			mockStorage: func(m *mocksWorker.MockRepo) {
				var lastRunAt *time.Time
				call := 0

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceInterval,
								Config:      []byte(`{"frequency":"daily", "interval":1}`),
								StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						expectedDue := time.Date(2026, 1, 1+call, 0, 0, 0, 0, time.UTC)

						assert.Equal(t, expectedDue, task.DueDate)

						call++
						lastRunAt = &next

						return nil
					}).Times(2)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
		{
			name:     "Generate first and second daily task when recurring task created in the past and never run",
			c:        clock.NewFake(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)),
			timeSkip: 0,
			mockStorage: func(m *mocksWorker.MockRepo) {
				var lastRunAt *time.Time
				call := 0

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						tasks := [][]taskdomain.RecurringTask{
							{},
							{{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceInterval,
								Config:      []byte(`{"frequency":"daily", "interval":1}`),
								StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
							}},
						}
						defer func() {
							call++
						}()
						return tasks[call], nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						expectedDue := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

						assert.Equal(t, expectedDue, task.DueDate)

						lastRunAt = &next

						return nil
					}).Times(1)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
		{
			name:     "Generate first daily task when recurring task created in the future",
			c:        clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
			timeSkip: 0,
			mockStorage: func(m *mocksWorker.MockRepo) {
				var lastRunAt *time.Time

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceInterval,
								Config:      []byte(`{"frequency":"daily", "interval":1}`),
								StartDate:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						expectedDue := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

						assert.Equal(t, expectedDue, task.DueDate)

						lastRunAt = &next

						return nil
					}).Times(1)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
		{
			name:     "Generate one daily task before end date",
			c:        clock.NewFake(time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)),
			timeSkip: 0,
			mockStorage: func(m *mocksWorker.MockRepo) {
				endDate := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)
				lastRunAt := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceInterval,
								Config:      []byte(`{"frequency":"daily", "interval":1}`),
								StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								EndDate:     &endDate,
								LastRunAt:   &lastRunAt,
								CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						expectedDue := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)

						assert.Equal(t, expectedDue, task.DueDate)

						lastRunAt = next

						return nil
					}).Times(1)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocksWorker.NewMockRepo(ctrl)
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			ctx := context.Background()
			validate := validator.New(validator.WithRequiredStructEnabled())

			schedulers := map[types.RecurrenceType]Scheduler{
				types.RecurrenceInterval: &IntervalScheduler{Validate: validate},
				types.RecurrenceOddDays:  &EvenOddDaysScheduler{Validate: validate},
				//types.RecurrenceYearlyOn: &YearlyDateScheduler{Validate: validate},
			}

			worker := New(mockRepo, tt.c, logger, schedulers)

			err := worker.RunOnce(ctx)
			if tt.want.firstRun != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.c.Add(tt.timeSkip)

			err = worker.RunOnce(ctx)
			if tt.want.secondRun != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWorker_OddDaysScheduler(t *testing.T) {
	type want struct {
		firstRun  error
		secondRun error
	}

	tests := []struct {
		name        string
		c           *clock.FakeClock
		timeSkip    time.Duration
		mockStorage func(m *mocksWorker.MockRepo)
		want        want
	}{
		{
			name:     "Generate first and second even day task. StartDate is even",
			c:        clock.NewFake(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)),
			timeSkip: 48 * time.Hour,
			mockStorage: func(m *mocksWorker.MockRepo) {
				var lastRunAt *time.Time
				call := 0

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceOddDays,
								Config:      []byte(`{"mode":"even"}`),
								StartDate:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						assert.Equal(t, 0, task.DueDate.Day()%2)

						call += 2
						lastRunAt = &next

						return nil
					}).Times(2)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
		{
			name:     "Generate first and second even day task. StartDate is odd",
			c:        clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
			timeSkip: 48 * time.Hour,
			mockStorage: func(m *mocksWorker.MockRepo) {
				var lastRunAt *time.Time
				call := 1

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						fmt.Println(lastRunAt)
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceOddDays,
								Config:      []byte(`{"mode":"even"}`),
								StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						assert.Equal(t, 0, task.DueDate.Day()%2)

						call += 2
						lastRunAt = &next

						return nil
					}).Times(2)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
		{
			name:     "Generate first and second odd day task. StartDate is odd",
			c:        clock.NewFake(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
			timeSkip: 48 * time.Hour,
			mockStorage: func(m *mocksWorker.MockRepo) {
				var lastRunAt *time.Time
				call := 0

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceOddDays,
								Config:      []byte(`{"mode":"odd"}`),
								StartDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						assert.Equal(t, 1, task.DueDate.Day()%2)

						call += 2
						lastRunAt = &next

						return nil
					}).Times(2)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
		{
			name:     "Generate first and second odd day task. StartDate is even",
			c:        clock.NewFake(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)),
			timeSkip: 48 * time.Hour,
			mockStorage: func(m *mocksWorker.MockRepo) {
				var lastRunAt *time.Time
				call := 1

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Type:        types.RecurrenceOddDays,
								Config:      []byte(`{"mode":"odd"}`),
								StartDate:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, types.StatusNew, task.Status)

						assert.Equal(t, 1, task.DueDate.Day()%2)

						call += 2
						lastRunAt = &next

						return nil
					}).Times(2)
			},
			want: want{
				firstRun:  nil,
				secondRun: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocksWorker.NewMockRepo(ctrl)
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			ctx := context.Background()
			validate := validator.New(validator.WithRequiredStructEnabled())

			schedulers := map[types.RecurrenceType]Scheduler{
				types.RecurrenceInterval: &IntervalScheduler{Validate: validate},
				types.RecurrenceOddDays:  &EvenOddDaysScheduler{Validate: validate},
				//types.RecurrenceYearlyOn: &YearlyDateScheduler{Validate: validate},
			}

			worker := New(mockRepo, tt.c, logger, schedulers)

			err := worker.RunOnce(ctx)
			if tt.want.firstRun != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.c.Add(tt.timeSkip)

			err = worker.RunOnce(ctx)
			if tt.want.secondRun != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
