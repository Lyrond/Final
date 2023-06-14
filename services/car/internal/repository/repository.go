package repository

import (
	"context"
	"errors"
	"final/services/car/internal/domain"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type carRep struct {
	db *pgxpool.Pool
}

type Car interface {
	Create(ctx context.Context, car *domain.Car) error
	GetByID(ctx context.Context, id int64) (*domain.Car, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, title string, brand []string) (interface{}, interface{})
}

func NewCarRep(db *pgxpool.Pool) *carRep {
	return &carRep{db: db}
}

func (s *carRep) Create(ctx context.Context, car *domain.Car) error {
	query := `
		INSERT INTO cars (car_title, year, brand)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []interface{}{car.Title, car.Year, car.Brand}

	return s.db.QueryRow(ctx, query, args...).Scan(&car.ID, &car.CreatedAt, &car.Version)
}

func (s *carRep) GetByID(ctx context.Context, id int64) (*domain.Car, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, year, brand, version
		FROM cars
		WHERE id = $1`

	var car domain.Car

	err := s.db.QueryRow(ctx, query, id).Scan(
		&car.ID,
		&car.CreatedAt,
		&car.Title,
		&car.Year,
		&car.Brand,
		&car.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &car, nil
}

func (s *carRep) GetAll(ctx context.Context, title string, brand []string) ([]*domain.Car, error) {
	query := fmt.Sprintf(`
		SELECT id, created_at, title, year, brand, version
		FROM cars
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (brand @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`)

	args := []any{title, brand}

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	defer rows.Close()

	cars := []*domain.Car{}

	for rows.Next() {
		var car domain.Car

		err := rows.Scan(
			&car.ID,
			&car.CreatedAt,
			&car.Title,
			&car.Year,
			&car.Brand,
			&car.Version,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	return cars, err
}

func (s *carRep) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM cars
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}
