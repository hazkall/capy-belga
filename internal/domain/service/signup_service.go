package service

import (
	"context"
	"fmt"

	"github.com/hazkall/capy-belga/internal/domain/entity"
	"github.com/hazkall/capy-belga/internal/domain/repository"
	"github.com/hazkall/capy-belga/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SignupService struct {
	Repo *repository.Repository
}

func (s *SignupService) UserClubStatus(ctx context.Context, signup *entity.SignupPayload) (bool, string, error) {

	cctx, span := telemetry.Tracer.Start(ctx, "SignupService.UserClubStatus",
		trace.WithAttributes(
			attribute.String("entity", "service"),
			attribute.String("email", signup.Email),
			attribute.String("club_name", signup.ClubName),
		),
	)

	defer span.End()

	userId, err := s.Repo.GetUserID(cctx, signup.Email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, "", err
	}

	active, plan, err := s.Repo.UserPlanStatus(cctx, userId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, "", err
	}

	return active, plan, nil
}

func (s *SignupService) SignupUser(ctx context.Context, signup *entity.SignupPayload) error {

	cctx, span := telemetry.Tracer.Start(ctx, "SignupService.SignupUser",
		trace.WithAttributes(
			attribute.String("entity", "service"),
			attribute.String("email", signup.Email),
			attribute.String("club_name", signup.ClubName),
		),
	)

	defer span.End()

	clubID, userID, err := s.Repo.GetUserIdClubID(cctx, signup.Email, signup.ClubName)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	active, err := s.Repo.UserState(cctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if !active {
		span.RecordError(fmt.Errorf("user %d is not active", userID))
		span.SetStatus(codes.Error, fmt.Errorf("user %d is not active", userID).Error())
		return fmt.Errorf("user %d is not active", userID)
	}
	return s.Repo.InsertUserClub(cctx, userID, clubID)
}

func (s *SignupService) CancelSignup(ctx context.Context, signup *entity.SignupPayload) error {

	cctx, span := telemetry.Tracer.Start(ctx, "SignupService.CancelSignup",
		trace.WithAttributes(
			attribute.String("entity", "service"),
			attribute.String("email", signup.Email),
		),
	)

	defer span.End()

	userID, err := s.Repo.GetUserID(cctx, signup.Email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	active, err := s.Repo.UserState(cctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if !active {
		span.RecordError(fmt.Errorf("user %d is not active", userID))
		span.SetStatus(codes.Error, fmt.Errorf("user %d is not active", userID).Error())
		return fmt.Errorf("user %d is not active", userID)
	}
	return s.Repo.CancelUserClub(cctx, userID)
}
