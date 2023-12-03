package main

import (
	"net/http"
)

func (app *Application) TerminalPage(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "terminal", nil); err != nil {
		app.ErrorLog.Println(err)
	}
}
