package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
)

func TerminalRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", controller.App.TerminalPage)
	mux.Get("/receipt", controller.App.Receipt)
	mux.Post("/payment-succeeded", controller.App.PaymentSucceeded)

	return mux
}
