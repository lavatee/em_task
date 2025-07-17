package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=1"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     string    `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name,omitempty"`
	Price       int     `json:"price,omitempty" binding:"min=0"`
	StartDate   string  `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	Subscription Subscription `json:"subscription"`
}

type SubscriptionListResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
