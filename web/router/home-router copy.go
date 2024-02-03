package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
	"github.com/nayeemnishaat/go-web-app/web/middleware"
)

func AdminRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Auth)

	mux.Get("/terminal", controller.App.TerminalPage)
	mux.Get("/all-sales", controller.App.AllSales)
	mux.Get("/all-subscriptions", controller.App.AllSubscriptions)

	return mux
}
