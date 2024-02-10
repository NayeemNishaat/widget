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
	mux.Get("/sales/{id}", controller.App.ShowSale)
	mux.Get("/subscriptions/{id}", controller.App.ShowSubscription)

	mux.Get("/all-users", controller.App.AllUsers)
	mux.Get("/all-users/{id}", controller.App.OneUser)

	return mux
}
