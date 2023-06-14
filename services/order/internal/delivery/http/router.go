package http

import (
	"final/services/order/internal/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type router struct {
	order OrderHandler
}

func NewRouter(orderService service.OrderService) *router {
	return &router{order: *NewHandler(orderService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/order/create", r.order.CreateOrder)
	router.HandlerFunc(http.MethodPost, "/v1/order/show", r.order.ShowOrder)

	return router
}
