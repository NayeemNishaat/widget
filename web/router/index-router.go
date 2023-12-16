package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/middleware"
)

func RootRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.SessionLoad)

	mux.Mount("/", HomeRouter())
	mux.Mount("/terminal", TerminalRouter())
	mux.Mount("/ecom", EcomRouter())

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
