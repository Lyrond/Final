package service

import (
	"context"
	"errors"
	"final/pkg/validator"
	"final/services/car/internal/domain"
	"final/services/car/internal/repository"
	"time"
)

var (
	ErrFailedValidation = errors.New("validation failed")
	ErrDuplicate        = errors.New("record duplication")
)

type CreateCarDTO struct {
	Title string   `json:"car_mark"`
	Year  int32    `json:"car_year"`
	Brand []string `json:"car_model"`
}

type CarService interface {
	CreateCar(ctx context.Context, car CreateCarDTO) error
	GetCarByID(ctx context.Context, id int64) (*domain.Car, error)
	GetCars(ctx context.Context, title string, brand []string) ([]*domain.Car, error)
}

type service struct {
	rep repository.Car
}

func New(repo repository.Car) *service {
	return &service{
		rep: repo,
	}
}

func (s *service) CreateCar(ctx context.Context, input CreateCarDTO) error {
	car := domain.Car{
		Title: input.Title,
		Year:  input.Year,
		Brand: input.Brand,
	}

	v := validator.New()

	if ValidateCar(v, &car); !v.Valid() {
		return ErrFailedValidation
	}

	err := s.rep.Create(ctx, &car)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetCarByID(ctx context.Context, id int64) (*domain.Car, error) {
	car, err := s.rep.GetByID(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return car, nil
}

func (s *service) GetCars(ctx context.Context, title string, brand []string) ([]*domain.Car, error) {

	var cars []*domain.Car

	cars, err := s.rep.GetAll(ctx, title, brand)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return cars, err
}

func ValidateCar(v *validator.Validator, car *domain.Car) {
	v.Check(car.Title != "", "title", "must be provided")
	v.Check(len(car.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(car.Year != 0, "year", "must be provided")
	v.Check(car.Year >= 1888, "year", "must be greater than 1900")
	v.Check(car.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(car.Brand != nil, "brand", "must be provided")
	v.Check(len(car.Brand) >= 1, "brand", "must contain at least 1 brand")
	v.Check(len(car.Brand) <= 5, "brand", "must not contain more than 5 brand")
	v.Check(validator.Unique(car.Brand), "brand", "must not contain duplicate values")
}
