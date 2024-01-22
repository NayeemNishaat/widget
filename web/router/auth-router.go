package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
)

func AuthRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/login", controller.App.LoginPage)

	return mux
}
