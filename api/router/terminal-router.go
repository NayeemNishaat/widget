package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
)

func UserRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/terminal", controller.App.TerminalPage)

	return mux
}
