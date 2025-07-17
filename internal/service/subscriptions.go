package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lavatee/subs/internal/model"
	"github.com/lavatee/subs/internal/repository"
	"github.com/sirupsen/logrus"
)

type SubscriptionsService struct {
	repo   *repository.Repository
	logger *logrus.Logger
}

func NewSubscriptionsService(repo *repository.Repository, logger *logrus.Logger) *SubscriptionsService {
	return &SubscriptionsService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SubscriptionsService) CreateSubscription(ctx context.Context, req model.CreateSubscriptionRequest) (model.Subscription, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		s.logger.Warnf("Invalid start date format: %v", err)
		return model.Subscription{}, err
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			s.logger.Warnf("Invalid end date format: %v", err)
			return model.Subscription{}, err
		}
		endDate = &parsedEndDate
	}

	subscription := model.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Subscriptions.CreateSubscription(ctx, subscription); err != nil {
		s.logger.Errorf("Failed to create subscription in repository: %v", err)
		return model.Subscription{}, err
	}

	return subscription, nil
}

func (s *SubscriptionsService) GetUserSubscriptions(ctx context.Context, userID uuid.UUID, serviceName string) ([]model.Subscription, error) {
	subscriptions, err := s.repo.Subscriptions.GetUserSubscriptions(ctx, userID, serviceName)
	if err != nil {
		s.logger.Errorf("Failed to get subscriptions from repository: %v", err)
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionsService) GetSubscription(ctx context.Context, id uuid.UUID) (model.Subscription, error) {
	subscription, err := s.repo.Subscriptions.GetSubscription(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get subscription from repository: %v", err)
		return model.Subscription{}, err
	}

	return subscription, nil
}

func (s *SubscriptionsService) UpdateSubscription(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) (model.Subscription, error) {
	existing, err := s.repo.Subscriptions.GetSubscription(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to get existing subscription for update: %v", err)
		return model.Subscription{}, err
	}

	if req.ServiceName != "" {
		existing.ServiceName = req.ServiceName
	}

	if req.Price != 0 {
		existing.Price = req.Price
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("01-2006", req.StartDate)
		if err != nil {
			s.logger.Warnf("Invalid start date format: %v", err)
			return model.Subscription{}, err
		}
		existing.StartDate = startDate
	}

	if req.EndDate != nil {
		if *req.EndDate == "" {
			existing.EndDate = nil
		} else {
			endDate, err := time.Parse("01-2006", *req.EndDate)
			if err != nil {
				s.logger.Warnf("Invalid end date format: %v", err)
				return model.Subscription{}, err
			}
			existing.EndDate = &endDate
		}
	}

	if err := s.repo.Subscriptions.UpdateSubscription(ctx, existing); err != nil {
		s.logger.Errorf("Failed to update subscription in repository: %v", err)
		return model.Subscription{}, err
	}

	return existing, nil
}

func (s *SubscriptionsService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Subscriptions.DeleteSubscription(ctx, id); err != nil {
		s.logger.Errorf("Failed to delete subscription from repository: %v", err)
		return err
	}

	return nil
}

func (s *SubscriptionsService) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startDate, endDate time.Time) (int, error) {
	total, err := s.repo.Subscriptions.GetTotalCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		s.logger.Errorf("Failed to calculate total cost in repository: %v", err)
		return 0, err
	}

	return total, nil
}
