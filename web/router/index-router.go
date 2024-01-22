package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/controller"
	"github.com/nayeemnishaat/go-web-app/web/middleware"
)

func RootRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.SessionLoad)

	mux.Mount("/", HomeRouter())

	mux.Mount("/auth", AuthRouter())

	mux.Mount("/ecom", EcomRouter())

	mux.Get("/terminal", controller.App.TerminalPage)

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
