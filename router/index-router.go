package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RootRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Mount("/", UserRouter())

	return mux
}
