package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type MiddlewareService struct {
	logger *slog.Logger
}

func NewMiddlewareService(logger *slog.Logger) *MiddlewareService {
	return &MiddlewareService{logger: logger}
}

func (s *MiddlewareService) LoggerMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		duration := time.Since(start)

		s.logger.Info("Request", "uri", r.RequestURI, "method", r.Method, "time", duration)
	})
}

func (s *MiddlewareService) RecoverPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				s.logger.Error("Request", "uri", r.RequestURI, "method", r.Method, "error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
