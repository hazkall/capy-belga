package worker

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/hazkall/capy-belga/internal/controller"
	"github.com/hazkall/capy-belga/internal/domain/entity"
	"github.com/hazkall/capy-belga/internal/domain/service"
	"github.com/hazkall/capy-belga/internal/mq"
	"github.com/hazkall/capy-belga/pkg/telemetry"
)

func StartPublishWorker(ctx context.Context, ch chan *controller.Message, m *mq.MQ) error {
	for club := range ch {
		if err := processClub(ctx, club, m); err != nil {
			return err
		}
	}
	return nil
}

func processClub(ctx context.Context, club *controller.Message, m *mq.MQ) error {

	_, span := telemetry.Tracer.Start(ctx, "processClubWorker",
		trace.WithAttributes(
			attribute.String("club_name", string(club.Data)),
			attribute.String("club_type", club.Type),
		),
	)

	defer span.End()

	var queueName string
	switch club.Type {
	case "create_discount_club":
		queueName = "discount_club_create"
	case "users":
		queueName = "users"
	case "discount_club_signup":
		queueName = "discount_club_signup"
	default:
		slog.Error("Unknown message type", "type", club.Type)
		return nil
	}

	j, err := json.Marshal(club)
	if err != nil {
		slog.Error("Error marshalling club entity", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = m.PublishMessage(j, queueName)

	if err != nil {
		slog.Error("Failed to publish message", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	slog.Info("Club created message published successfully", "club", club)

	return nil
}

func ConsumeCreateClub(ctx context.Context, m *mq.MQ, clubService *service.ClubService) error {

	club := new(entity.Club)
	msg := new(controller.Message)

	deliveries, err := m.ConsumeMessages("discount_club_create")
	if err != nil {
		slog.Error("Failed to consume messages", "error", err)
		return err
	}

	for d := range deliveries {
		cctx, span := telemetry.Tracer.Start(ctx, "ConsumeCreateClubWorker",
			trace.WithAttributes(
				attribute.String("entity", "worker"),
				attribute.String("queue_name", "discount_club_create"),
				attribute.String("message_type", d.Type),
			),
		)
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			slog.Error("Error unmarshalling club entity", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		if err := json.Unmarshal(msg.Data, &club); err != nil {
			slog.Error("Error unmarshalling club entity data", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		err := clubService.CreateClub(cctx, club)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				slog.Warn("Duplicate club detected, skipping", "name", club.Name)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				d.Ack(false)
				continue
			}
			slog.Error("Error inserting club into database", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		span.SetAttributes(
			attribute.String("club_name", club.Name),
			attribute.String("club_type", club.PlanType),
			attribute.String("AcquisitionChannel", club.AquisitionChannel),
			attribute.String("AcquisitionLocation", club.AquisitionLocation),
		)
		span.End()
		d.Ack(false)
	}

	return nil

}

func ConsumeUser(ctx context.Context, m *mq.MQ, userService *service.UserService) error {
	user := new(entity.User)
	msg := new(controller.Message)

	deliveries, err := m.ConsumeMessages("users")
	if err != nil {
		slog.Error("Failed to consume messages", "error", err)
		return err
	}

	for d := range deliveries {
		cctx, span := telemetry.Tracer.Start(ctx, "ConsumeUserWorker",
			trace.WithAttributes(
				attribute.String("entity", "worker"),
				attribute.String("queue_name", "users"),
				attribute.String("message_type", d.Type),
			),
		)
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			slog.Error("Error unmarshalling user entity", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		if err := json.Unmarshal(msg.Data, &user); err != nil {
			slog.Error("Error unmarshalling user entity data", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		err := userService.CreateUser(cctx, user)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				slog.Warn("Duplicate user detected, skipping", "email", user.Email)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				d.Ack(false)
				continue
			}
			slog.Error("Error inserting user into database", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}
		span.SetAttributes(
			attribute.String("email", user.Email),
		)
		span.End()
		d.Ack(false)
	}
	return nil
}

func ConsumeClubSignup(ctx context.Context, m *mq.MQ, signupService *service.SignupService) error {
	signup := new(entity.SignupPayload)
	msg := new(controller.Message)

	deliveries, err := m.ConsumeMessages("discount_club_signup")
	if err != nil {
		slog.Error("Failed to consume messages", "error", err)

		return err
	}

	for d := range deliveries {
		cctx, span := telemetry.Tracer.Start(ctx, "ConsumeClubSignupWorker",
			trace.WithAttributes(
				attribute.String("entity", "worker"),
				attribute.String("queue_name", "discount_club_signup"),
				attribute.String("message_type", d.Type),
			),
		)
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			slog.Error("Error unmarshalling signup entity", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		if err := json.Unmarshal(msg.Data, &signup); err != nil {
			slog.Error("Error unmarshalling signup entity data", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		err := signupService.SignupUser(cctx, signup)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				slog.Warn("Duplicate signup detected, skipping", "email", signup.Email, "club", signup.ClubName)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				d.Ack(false)
				continue
			}
			slog.Error("Error processing club signup", "error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			d.Nack(false, true)
			continue
		}

		span.SetAttributes(
			attribute.String("email", signup.Email),
			attribute.String("club_name", signup.ClubName),
		)
		span.End()

		telemetry.PlanGauge.Record(
			cctx,
			1,
			metric.WithAttributes(
				attribute.String("email", signup.Email),
				attribute.String("club_name", signup.ClubName),
				attribute.String("action", "signup"),
			),
		)

		telemetry.NewPlanCounter.Add(
			cctx,
			1,
			metric.WithAttributes(
				attribute.String("email", signup.Email),
				attribute.String("club_name", signup.ClubName),
				attribute.String("action", "signup"),
			),
		)

		d.Ack(false)
	}

	return nil
}
