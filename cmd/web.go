package cmd

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	httpdelivery "github.com/pevin/pevin-golang-training-beginner/delivery/http"
)

type WebCommand struct{}

func (c *WebCommand) Execute() {
	paymentcodeHandler := httpdelivery.NewPaymentCodeRouteHandler(initPaymentUsecase())
	http.HandleFunc("/payment-codes", paymentcodeHandler.PaymentCodeRouteHandler)
	http.HandleFunc("/payment-codes/", paymentcodeHandler.PaymentCodeRouteHandler)

	http.HandleFunc("/health", httpdelivery.HealthHandler)
	http.HandleFunc("/hello-world", httpdelivery.HelloWorldHandler)

	http.HandleFunc("/", httpdelivery.NotFoundHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
