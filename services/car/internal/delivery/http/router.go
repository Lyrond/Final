package http

import (
	"final/services/car/internal/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type router struct {
	car CarHandler
}

func NewRouter(carService service.CarService) *router {
	return &router{car: *NewHandler(carService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/cars", r.car.CreateCarHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cars/:id", r.car.ShowCarHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cars", r.car.ListCarsHandler)

	return router
}
