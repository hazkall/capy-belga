package service

import (
	"context"

	"github.com/hazkall/capy-belga/internal/domain/entity"
	"github.com/hazkall/capy-belga/internal/domain/repository"
	"github.com/hazkall/capy-belga/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ClubService struct {
	Repo *repository.Repository
}

func (s *ClubService) CreateClub(ctx context.Context, club *entity.Club) error {

	cctx, span := telemetry.Tracer.Start(ctx, "ClubService.CreateClub",
		trace.WithAttributes(
			attribute.String("club_name", club.Name),
			attribute.String("club_type", club.PlanType),
			attribute.String("AcquisitionChannel", club.AquisitionChannel),
			attribute.String("AcquisitionLocation", club.AquisitionLocation),
		),
	)
	defer span.End()

	return s.Repo.InsertClub(cctx, club)
}
