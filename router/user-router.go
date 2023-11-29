package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/controller"
)

func UserRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/user", controller.App.User)

	return mux
}
