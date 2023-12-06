package controller

import (
	"net/http"

	"github.com/nayeemnishaat/go-web-app/web/template"
)

func (app *Application) ChargeOncePage(w http.ResponseWriter, r *http.Request) {
	widget := struct {
		ID             int
		Name           string
		Description    string
		InventoryLevel int
		Price          int
	}{ID: 1, Name: "Custom Widget", Description: "Amazing", InventoryLevel: 10, Price: 1000}

	if err := app.RenderTemplate(w, r, "buy-once", &template.TemplateData{StringMap: map[string]string{"publishable_key": app.Stripe.Key}, Data: map[string]any{"widget": widget}}, "stripe-js"); err != nil {
		app.ErrorLog.Println(err)
	}
}
