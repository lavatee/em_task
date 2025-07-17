package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lavatee/subs/internal/model"
)

type Subscriptions interface {
	CreateSubscription(ctx context.Context, sub model.Subscription) error
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID, serviceName string) ([]model.Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub model.Subscription) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startDate, endDate time.Time) (int, error)
}

type Repository struct {
	Subscriptions
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Subscriptions: NewSubscriptionsPostgres(db),
	}
}
