package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
)

func EcomRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/bronze", controller.App.BronzePlan)
	mux.Get("/widget/{id}", controller.App.ChargeOncePage)

	mux.Get("/receipt", controller.App.Receipt)
	mux.Get("/virtual-receipt", controller.App.VirtualReceipt)
	mux.Post("/payment-succeeded", controller.App.PaymentSucceeded)
	mux.Post("/virtual-payment-succeeded", controller.App.VirtualPaymentSucceeded)

	return mux
}
