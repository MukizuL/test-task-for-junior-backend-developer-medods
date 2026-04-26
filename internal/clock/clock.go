package clock

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

type FakeClock struct {
	now time.Time
}

func (f *FakeClock) Now() time.Time {
	return f.now
}

func (f *FakeClock) Add(d time.Duration) {
	f.now = f.now.Add(d)
}

func NewFake(date time.Time) *FakeClock {
	return &FakeClock{
		now: date.UTC(),
	}
}
