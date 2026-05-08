package task

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	repo     Repository
	validate *validator.Validate
	now      func() time.Time
}

func NewService(repo Repository, validate *validator.Validate) *Service {
	return &Service{
		repo:     repo,
		validate: validate,
		now:      func() time.Time { return time.Now().UTC() },
	}
}
