package service

import (
	"context"

	"github.com/hazkall/capy-belga/internal/domain/entity"
	"github.com/hazkall/capy-belga/internal/domain/repository"
	"github.com/hazkall/capy-belga/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UserService struct {
	Repo *repository.Repository
}

func (s *UserService) CreateUser(ctx context.Context, user *entity.User) error {

	cctx, span := telemetry.Tracer.Start(ctx, "UserService.CreateUser",
		trace.WithAttributes(
			attribute.String("entity", "service"),
			attribute.String("email", user.Email),
		),
	)
	defer span.End()

	return s.Repo.InsertUser(cctx, user)
}

func (s *UserService) UserState(ctx context.Context, email string) (bool, error) {

	cctx, span := telemetry.Tracer.Start(ctx, "UserService.UserState",
		trace.WithAttributes(
			attribute.String("entity", "service"),
			attribute.String("email", email),
		),
	)
	defer span.End()

	userId, err := s.Repo.GetUserID(cctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	status, err := s.Repo.UserState(cctx, userId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	return status, nil
}
