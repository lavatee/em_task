package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lavatee/subs/internal/model"
	"github.com/lavatee/subs/internal/repository"
	"github.com/sirupsen/logrus"
)

type Subscriptions interface {
	CreateSubscription(ctx context.Context, request model.CreateSubscriptionRequest) (model.Subscription, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID, serviceName string) ([]model.Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (model.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, request model.UpdateSubscriptionRequest) (model.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startDate, endDate time.Time) (int, error)
}

type Service struct {
	Subscriptions
}

func NewService(repo *repository.Repository, logger *logrus.Logger) *Service {
	return &Service{
		Subscriptions: NewSubscriptionsService(repo, logger),
	}
}
