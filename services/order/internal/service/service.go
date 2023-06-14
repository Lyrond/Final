package service

import (
	"context"
	"errors"
	"final/pkg/validator"
	"final/services/order/internal/domain"
	"final/services/order/internal/repository"
)

var (
	ErrFailedValidation = errors.New("validation failed")
)

type CreateOrderDTO struct {
	CarID int64  `json:"car_id"`
	Email string `json:"email"`
}

type OrderService interface {
	Create(ctx context.Context, order CreateOrderDTO) error
	Show(ctx context.Context, email string) ([]*domain.Order, error)
}

type service struct {
	repo repository.Order
}

func New(repo repository.Order) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, input CreateOrderDTO) error {
	order := domain.Order{
		CarID: input.CarID,
		Email: input.Email,
	}
	v := validator.New()
	validateEmail(v, order.Email)
	if !v.Valid() {
		return ErrFailedValidation
	}
	//TODO: add validator for CarID and email
	err := s.repo.Insert(ctx, &order)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) Show(ctx context.Context, email string) ([]*domain.Order, error) {
	v := validator.New()
	validateEmail(v, email)
	if !v.Valid() {
		return nil, ErrFailedValidation
	}

	orders, err := s.repo.GetByEmail(ctx, &email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return orders, nil
}

func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
