package main

import (
	"net/http"

	"github.com/nayeemnishaat/go-web-app/api/lib"
)

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.authenticateToken(r)
		if err != nil {
			lib.InvalidCredentials(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
