package handler

import (
	"final/services/payment/internal/delivery/http"
	"github.com/gorilla/mux"
)

func NewRouter(paymentHandler *handler.PaymentHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/process-payment", paymentHandler.ProcessPayment).Methods("POST")
	return r
}
