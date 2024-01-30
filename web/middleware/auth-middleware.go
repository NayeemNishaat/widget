package middleware

import (
	"net/http"

	"github.com/nayeemnishaat/go-web-app/web/lib"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app := lib.GetConfig()

		if !app.Session.Exists(r.Context(), "userID") {
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
