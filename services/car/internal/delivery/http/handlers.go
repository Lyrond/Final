package http

import (
	"errors"
	"final/pkg/request"
	"final/services/car/internal/repository"
	"final/services/car/internal/service"
	"net/http"
)

type CarHandler struct {
	carService service.CarService
}

func NewHandler(service service.CarService) *CarHandler {
	return &CarHandler{carService: service}
}

func (h *CarHandler) CreateCarHandler(w http.ResponseWriter, r *http.Request) {
	var input service.CreateCarDTO

	err := request.ReadJSON(w, r, &input)
	if err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	car := service.CreateCarDTO{
		Title: input.Title,
		Year:  input.Year,
		Brand: input.Brand,
	}

	err = h.carService.CreateCar(r.Context(), car)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			request.BadRequestResponse(w, r, err)
			return
		case errors.Is(err, service.ErrDuplicate):
			request.RecordDuplicationResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}

	err = request.WriteJSON(w, http.StatusCreated, map[string]any{"car": car}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}

}

func (h *CarHandler) ShowCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := request.ReadIDParam(r)
	if err != nil {
		request.NotFoundResponse(w, r)
		return
	}

	car, err := h.carService.GetCarByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"car": car}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *CarHandler) ListCarsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Brand []string
	}
	qs := r.URL.Query()
	input.Title = request.ReadString(qs, "title", "")
	input.Brand = request.ReadCSV(qs, "brand", []string{})

	cars, err := h.carService.GetCars(r.Context(), input.Title, input.Brand)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"cars": cars}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}
