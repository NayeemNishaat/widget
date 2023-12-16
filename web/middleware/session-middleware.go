package middleware

import (
	"net/http"

	"github.com/nayeemnishaat/go-web-app/web/lib"
)

func SessionLoad(next http.Handler) http.Handler {
	return lib.Session.LoadAndSave(next)
}
