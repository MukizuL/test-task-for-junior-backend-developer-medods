package worker

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	mocksWorker "example.com/taskservice/cmd/worker/mocks"
	"example.com/taskservice/internal/clock"
	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorker_GeneratesTask(t *testing.T) {
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
			name:     "Generate first and second daily task",
			c:        clock.NewFake(time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC)),
			timeSkip: 24 * time.Hour,
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
								StartDate:   time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().CreateAndUpdateLastRunAt(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task, next time.Time) error {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, taskdomain.StatusNew, task.Status)

						expectedDue := time.Date(2026, 4, 26+call, 10, 0, 0, 0, time.UTC)

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
				types.RecurrenceOddDays:  &OddDaysScheduler{Validate: validate},
				types.RecurrenceEvenDays: &EvenDaysScheduler{Validate: validate},
				types.RecurrenceYearlyOn: &YearlyDateScheduler{Validate: validate},
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
