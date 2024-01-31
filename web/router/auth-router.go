package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
)

func AuthRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/login", controller.App.LoginPage)
	mux.Post("/login", controller.App.PostLoginPage)
	mux.Get("/logout", controller.App.Logout)
	mux.Get("/forgot-password", controller.App.ForgotPassword)

	return mux
}
