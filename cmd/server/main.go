package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/hazkall/capy-belga/internal/controller"
	"github.com/hazkall/capy-belga/internal/domain/repository"
	"github.com/hazkall/capy-belga/internal/domain/service"
	"github.com/hazkall/capy-belga/internal/mq"
	"github.com/hazkall/capy-belga/internal/router"
	"github.com/hazkall/capy-belga/internal/server"
	"github.com/hazkall/capy-belga/internal/worker"
	"github.com/hazkall/capy-belga/pkg/logger"
	"github.com/hazkall/capy-belga/pkg/telemetry"
)

func main() {

	logger.Start("true")
	ctx := context.Background()
	clubChannel := make(chan *controller.Message, 100)

	slog.Info("Starting OpenTelemetry Tracing")

	t := telemetry.TraceInit(ctx, "capy-belga-tracer")

	defer func() {
		if err := t.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down exporter", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("Starting OpenTelemetry Metrics")
	me, exp := telemetry.MetricInit(ctx, "capy-belga-metrics")

	slog.Info("Starting OpenTelemetry Go Runtime Metrics")
	telemetry.RuntimeStart(me)

	slog.Info("Starting OpenTelemetry Metrics")
	if err := telemetry.MetricsStart(); err != nil {
		slog.Error("Error starting metrics", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := exp.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down exporter", "error", err)
			os.Exit(1)
		}
	}()

	m, err := mq.NewMQ(os.Getenv("RABBITMQ_URL"))

	if err != nil {
		slog.Error("Failed to create message queue", "error", err)
		return
	}

	defer m.Close()

	queueNames := []string{"discount_club_create", "users", "discount_club_signup"}

	if err := m.DeclareQueues(queueNames); err != nil {
		slog.Error("Failed to declare queues", "error", err)
		return
	}

	serverAddress := ":8080"
	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second
	idleTimeout := 15 * time.Second

	slog.Info("Handlers pipeline initialized")

	repo := repository.NewRepository()
	clubService := service.ClubService{Repo: repo}
	userService := service.UserService{Repo: repo}
	signupService := service.SignupService{Repo: repo}

	deps := &router.HandlerDeps{
		ClubChannel:   clubChannel,
		UserService:   &userService,
		ClubService:   &clubService,
		SignupService: &signupService,
	}

	router.HandlersPipeline(deps)

	for i := 0; i < 5; i++ {
		go func(workerID int) {
			slog.Info("Starting worker to process clubs", "worker_id", workerID)
			if err := worker.StartPublishWorker(ctx, clubChannel, m); err != nil {
				slog.Error("Failed to start publishing worker", "error", err)
				return
			}
		}(i + 1)
	}

	go func() {
		slog.Info("Starting worker to consume clubs")
		if err := worker.ConsumeCreateClub(ctx, m, &clubService); err != nil {
			slog.Error("Failed to start consuming worker", "error", err)
			return
		}
	}()

	go func() {
		slog.Info("Starting worker to consume users")
		if err := worker.ConsumeUser(ctx, m, &userService); err != nil {
			slog.Error("Failed to start consuming worker", "error", err)
			return
		}
	}()

	go func() {
		slog.Info("Starting worker to consume discount club signups")
		if err := worker.ConsumeClubSignup(ctx, m, &signupService); err != nil {
			slog.Error("Failed to start consuming worker", "error", err)
			return
		}
	}()

	server.StartServer(ctx, serverAddress, readTimeout, writeTimeout, idleTimeout)
}
