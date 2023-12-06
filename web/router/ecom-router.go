package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
)

func EcomRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/charge-once", controller.App.ChargeOncePage)

	return mux
}
