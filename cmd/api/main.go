package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/taskservice/internal/clock"
	infrastructurepostgres "example.com/taskservice/internal/infrastructure/postgres"
	postgresrepo "example.com/taskservice/internal/repository/postgres"
	transporthttp "example.com/taskservice/internal/transport/http"
	swaggerdocs "example.com/taskservice/internal/transport/http/docs"
	httphandlers "example.com/taskservice/internal/transport/http/handlers"
	"example.com/taskservice/internal/transport/http/middleware"
	"example.com/taskservice/internal/types"
	"example.com/taskservice/internal/usecase/task"
	"example.com/taskservice/internal/worker"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	err := loadConfig()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		os.Exit(0)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	pool, err := infrastructurepostgres.Open(ctx, viper.GetString("DATABASE_DSN"))
	if err != nil {
		logger.Error("open postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	validate := validator.New(validator.WithRequiredStructEnabled())

	schedulers := map[types.RecurrenceType]worker.Scheduler{
		types.RecurrenceInterval:      &worker.IntervalScheduler{Validate: validate},
		types.RecurrenceEvenOddDays:   &worker.EvenOddDaysScheduler{Validate: validate},
		types.RecurrenceSpecificDates: &worker.SpecificDateScheduler{Validate: validate},
	}

	taskRepo := postgresrepo.New(pool)
	taskUsecase := task.NewService(taskRepo, validate, schedulers)
	taskHandler := httphandlers.NewTaskHandler(taskUsecase)
	docsHandler := swaggerdocs.NewHandler()
	mw := middleware.NewMiddlewareService(logger)
	router := transporthttp.NewRouter(taskHandler, docsHandler, mw.RecoverPanic, mw.LoggerMW)
	workerRecTask := worker.New(taskRepo, clock.RealClock{}, logger, schedulers)

	g.Go(func() error {
		return workerRecTask.Run(gCtx, viper.GetDuration("WORKER_TICK"))
	})

	server := &http.Server{
		Addr:              viper.GetString("HTTP_ADDR"),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-gCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown http server", "error", err)
		}
	}()

	logger.Info("http server started", "addr", viper.GetString("HTTP_ADDR"))

	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.ErrorContext(gCtx, "Emergency stop...", "error", err)
	}
}

func loadConfig() error {
	viper.SetDefault("HTTP_ADDR", ":8080")
	viper.SetDefault("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/taskservice?sslmode=disable")
	viper.SetDefault("WORKER_TICK", "5s")

	err := viper.BindEnv("HTTP_ADDR")
	if err != nil {
		return err
	}

	err = viper.BindEnv("DATABASE_DSN")
	if err != nil {
		return err
	}

	err = viper.BindEnv("WORKER_TICK")
	if err != nil {
		return err
	}

	return nil
}
