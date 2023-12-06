package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RootRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Mount("/terminal", TerminalRouter())
	mux.Mount("/ecom", EcomRouter())

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
