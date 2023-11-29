package controller

import (
	"net/http"
)

func (app *Application) User(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("OK")
}
