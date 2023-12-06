package controller

import (
	"net/http"
)

func (app *Application) ChargeOncePage(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "buy-once", nil, "stripe-js"); err != nil {
		app.ErrorLog.Println(err)
	}
}
