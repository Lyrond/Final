package repository

import (
	"context"
	"errors"
	"final/services/order/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type orderRepo struct {
	db *pgxpool.Pool
}

type Order interface {
	Insert(ctx context.Context, order *domain.Order) error
	GetByEmail(ctx context.Context, email *string) ([]*domain.Order, error)
}

func NewOrderRepo(db *pgxpool.Pool) *orderRepo {
	return &orderRepo{db: db}
}

func (s *orderRepo) Insert(ctx context.Context, order *domain.Order) error {
	query := `
	INSERT INTO orders (car_mark, email)
	VALUES ($1, $2)
	RETURNING id, created_at`

	args := []any{order.CarID, order.Email}

	err := s.db.QueryRow(ctx, query, args...).Scan(&order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *orderRepo) GetByEmail(ctx context.Context, email *string) ([]*domain.Order, error) {
	query := `
	SELECT car_id, car_mark, email
	FROM orders
	WHERE email = $1`
	rows, err := s.db.Query(ctx, query, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.CarID,
			&order.Email,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}
