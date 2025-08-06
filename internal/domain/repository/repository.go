package repository

import (
	"context"
	"log/slog"
	"os"

	"github.com/hazkall/capy-belga/internal/db"
	"github.com/hazkall/capy-belga/internal/domain/entity"
	"github.com/hazkall/capy-belga/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	db *db.Postgres
}

func NewRepository() *Repository {
	dbConn, err := db.NewPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		slog.Error("Failed to create database connection", "error", err)
		return nil
	}
	return &Repository{db: dbConn}
}

func (r *Repository) GetUserIdClubID(ctx context.Context, email, clubName string) (clubId, userId int64, err error) {

	_, span := telemetry.Tracer.Start(ctx, "Repository.GetUserIdClubID",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.String("email", email),
			attribute.String("club_name", clubName),
		),
	)

	defer span.End()

	query := `
		SELECT u.id, c.id
		FROM users u
		JOIN clubs c ON 1=1
		WHERE u.email = $1 AND c.name = $2
	`

	span.SetAttributes(attribute.String("db.statement", query))

	row := r.db.DB.QueryRow(query, email, clubName)
	err = row.Scan(&userId, &clubId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return
}

func (r *Repository) GetUserID(ctx context.Context, email string) (int64, error) {
	_, span := telemetry.Tracer.Start(ctx, "Repository.GetUserID",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.String("email", email),
		),
	)
	defer span.End()

	query := `SELECT id FROM users WHERE email = $1`

	span.SetAttributes(attribute.String("db.statement", query))

	var userID int64
	err := r.db.DB.QueryRow(query, email).Scan(&userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	return userID, nil
}

func (r *Repository) GetClubID(ctx context.Context, name string) (int64, error) {
	_, span := telemetry.Tracer.Start(ctx, "Repository.GetClubID",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.String("name", name),
		),
	)
	defer span.End()

	query := `
		SELECT id
		FROM clubs
		WHERE name = $1
	`

	span.SetAttributes(attribute.String("db.statement", query))

	var clubID int64
	err := r.db.DB.QueryRow(query, name).Scan(&clubID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	return clubID, nil
}

func (r *Repository) UserState(ctx context.Context, userID int64) (active bool, err error) {
	_, span := telemetry.Tracer.Start(ctx, "Repository.UserState",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.Int64("user_id", userID),
		),
	)
	defer span.End()

	query := `
		SELECT active
		FROM users
		WHERE id = $1
	`
	span.SetAttributes(attribute.String("db.statement", query))

	err = r.db.DB.QueryRow(query, userID).Scan(&active)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}
	return active, nil
}

func (r *Repository) UserPlanStatus(ctx context.Context, userID int64) (active bool, planType string, err error) {
	_, span := telemetry.Tracer.Start(ctx, "Repository.UserPlanStatus",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.Int64("user_id", userID),
		),
	)
	defer span.End()

	query := `
		SELECT uc.active, c.plan_type
		FROM user_club uc
		JOIN clubs c ON uc.club_id = c.id
		WHERE uc.user_id = $1
		LIMIT 1
	`
	span.SetAttributes(attribute.String("db.statement", query))

	err = r.db.DB.QueryRow(query, userID).Scan(&active, &planType)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return
}

func (r *Repository) InsertUserClub(ctx context.Context, userID int64, clubID int64) error {
	_, span := telemetry.Tracer.Start(ctx, "Repository.InsertUserClub",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.Int64("user_id", userID),
			attribute.Int64("club_id", clubID),
		),
	)
	defer span.End()

	query := `
		INSERT INTO user_club (user_id, club_id)
		VALUES ($1, $2)
	`
	span.SetAttributes(attribute.String("db.statement", query))

	_, err := r.db.DB.Exec(query, userID, clubID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (r *Repository) InsertUser(ctx context.Context, user *entity.User) error {
	_, span := telemetry.Tracer.Start(ctx, "Repository.InsertUser",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.String("email", user.Email),
		),
	)
	defer span.End()

	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
	`

	span.SetAttributes(attribute.String("db.statement", query))

	_, err := r.db.DB.Exec(query, user.Name, user.Email)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

func (r *Repository) InsertClub(ctx context.Context, club *entity.Club) error {
	_, span := telemetry.Tracer.Start(ctx, "Repository.InsertClub",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.String("name", club.Name),
		),
	)
	defer span.End()

	query := `
		INSERT INTO clubs (name, description, aquisition_channel, aquisition_location, plan_type)
		VALUES ($1, $2, $3, $4, $5)
	`
	span.SetAttributes(attribute.String("db.statement", query))

	_, err := r.db.DB.Exec(query, club.Name, club.Description, club.AquisitionChannel, club.AquisitionLocation, club.PlanType)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

func (r *Repository) CancelUserClub(ctx context.Context, userID int64) error {
	_, span := telemetry.Tracer.Start(ctx, "Repository.CancelUserClub",
		trace.WithAttributes(
			attribute.String("entity", "repository"),
			attribute.Int64("user_id", userID),
		),
	)
	defer span.End()

	query := `
		UPDATE user_club
		SET active = false
		WHERE user_id = $1
	`
	span.SetAttributes(attribute.String("db.statement", query))
	_, err := r.db.DB.Exec(query, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
