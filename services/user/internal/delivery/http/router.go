package http

import (
	"final/services/user/internal/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type router struct {
	user UserHandler
}

func NewRouter(userService service.UserService) *router {
	return &router{user: *NewHandler(userService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/user/signup", r.user.RegisterUser)
	router.HandlerFunc(http.MethodPost, "/v1/user/signin", r.user.LoginUser)

	return router
}
