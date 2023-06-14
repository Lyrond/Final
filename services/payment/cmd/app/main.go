package app

import (
	"final/services/payment/internal/delivery/http"
	"final/services/payment/internal/service"
	"log"
	"net/http"
)

func main() {
	paymentService := service.NewPaymentService()

	paymentHandler := handler.NewPaymentHandler(paymentService)

	r := router.NewRouter(paymentHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("service started on :8080")
	log.Fatal(server.ListenAndServe())
}
