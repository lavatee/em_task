package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lavatee/subs/internal/model"
)

type SubscriptionsPostgres struct {
	db *sqlx.DB
}

func NewSubscriptionsPostgres(db *sqlx.DB) *SubscriptionsPostgres {
	return &SubscriptionsPostgres{
		db: db,
	}
}

func (r *SubscriptionsPostgres) CreateSubscription(ctx context.Context, sub model.Subscription) error {
	query := fmt.Sprintf(`INSERT INTO %s
	(id, service_name, price, user_id, start_date, end_date, created_at)
	VALUES (:id, :service_name, :price, :user_id, :start_date, :end_date, :created_at)`, subscriptionsTable)
	_, err := r.db.NamedExecContext(ctx, query, sub)
	return err
}

func (r *SubscriptionsPostgres) GetUserSubscriptions(ctx context.Context, userID uuid.UUID, serviceName string) ([]model.Subscription, error) {
	var userIDArg interface{} = userID
	if userID == uuid.Nil {
		userIDArg = nil
	}

	var serviceNameArg interface{} = serviceName
	if serviceName == "" {
		serviceNameArg = nil
	}

	query := fmt.Sprintf(`SELECT id, service_name, price, user_id, start_date, end_date, created_at
    FROM %s
    WHERE ($1::uuid IS NULL OR user_id = $1)
    AND ($2::text IS NULL OR service_name = $2)
    ORDER BY created_at DESC`, subscriptionsTable)

	var subs []model.Subscription
	if err := r.db.SelectContext(ctx, &subs, query, userIDArg, serviceNameArg); err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	return subs, nil
}

func (r *SubscriptionsPostgres) GetSubscription(ctx context.Context, id uuid.UUID) (model.Subscription, error) {
	query := fmt.Sprintf(`SELECT id, service_name, price, user_id, start_date, end_date, created_at
	FROM %s
	WHERE id = $1`, subscriptionsTable)
	var sub model.Subscription
	if err := r.db.GetContext(ctx, &sub, query, id); err != nil {
		return model.Subscription{}, err
	}
	return sub, nil
}

func (r *SubscriptionsPostgres) UpdateSubscription(ctx context.Context, sub model.Subscription) error {
	query := fmt.Sprintf(`UPDATE %s
	SET service_name = :service_name,
	price = :price,
	start_date = :start_date,
	end_date = :end_date
	WHERE id = :id`, subscriptionsTable)
	_, err := r.db.NamedExecContext(ctx, query, sub)
	return err
}

func (r *SubscriptionsPostgres) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, subscriptionsTable)
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("Subscription with this id does not exist")
	}
	return nil
}

func (r *SubscriptionsPostgres) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startDate, endDate time.Time) (int, error) {
	var userIDArg interface{} = userID
	if userID == uuid.Nil {
		userIDArg = nil
	}

	var serviceNameArg interface{} = serviceName
	if serviceName == "" {
		serviceNameArg = nil
	}

	var startDateArg interface{} = startDate
	if startDate.IsZero() {
		startDateArg = nil
	}

	var endDateArg interface{} = endDate
	if endDate.IsZero() {
		endDateArg = nil
	}

	query := fmt.Sprintf(`SELECT COALESCE(SUM(price), 0) 
    FROM %s 
    WHERE ($1::uuid IS NULL OR user_id = $1)
    AND ($2::text IS NULL OR service_name = $2)
    AND ($3::timestamp IS NULL OR (end_date IS NULL OR end_date >= $3))
    AND ($4::timestamp IS NULL OR start_date <= $4)`, subscriptionsTable)

	var total int
	if err := r.db.GetContext(ctx, &total, query, userIDArg, serviceNameArg, startDateArg, endDateArg); err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	return total, nil
}
