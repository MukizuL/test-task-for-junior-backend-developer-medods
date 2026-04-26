package worker

import (
	"context"
	"testing"
	"time"

	"example.com/taskservice/internal/clock"
	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/usecase/task/mocks"
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
		mockStorage func(m *mocks.MockRepo)
		want        want
	}{
		{
			name:     "Generate first and second daily task",
			c:        clock.NewFake(time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC)),
			timeSkip: 24 * time.Hour,
			mockStorage: func(m *mocks.MockRepo) {
				var lastRunAt *time.Time
				call := 0

				m.EXPECT().GetDueRecurringTasks(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]taskdomain.RecurringTask, error) {
						return []taskdomain.RecurringTask{
							{
								ID:          1,
								Title:       "Task 1",
								Description: "Description 1",
								Frequency:   "daily",
								Interval:    1,
								StartDate:   time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC),
								EndDate:     nil,
								LastRunAt:   lastRunAt,
								CreatedAt:   time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC),
								UpdatedAt:   time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC),
							},
						}, nil
					}).Times(2)

				m.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&taskdomain.Task{})).
					DoAndReturn(func(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
						assert.Equal(t, int64(1), *task.RecurringTaskID)
						assert.Equal(t, "Task 1", task.Title)
						assert.Equal(t, "Description 1", task.Description)
						assert.Equal(t, taskdomain.StatusNew, task.Status)

						expectedDue := time.Date(2026, 4, 26+call, 10, 0, 0, 0, time.UTC)

						assert.Equal(t, expectedDue, task.DueDate)

						call++

						return task, nil
					}).Times(2)

				m.EXPECT().UpdateLastRunAt(gomock.Any(), int64(1), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int64, next time.Time) error {
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

			mockRepo := mocks.NewMockRepo(ctrl)
			if tt.mockStorage != nil {
				tt.mockStorage(mockRepo)
			}

			worker := New(context.Background(), mockRepo, tt.c)

			err := worker.RunOnce()
			if tt.want.firstRun != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.c.Add(tt.timeSkip)

			err = worker.RunOnce()
			if tt.want.secondRun != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
