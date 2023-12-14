package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/web/template"
)

func (app *Application) ChargeOncePage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	if err := app.RenderTemplate(w, r, "buy-once", &template.TemplateData{StringMap: map[string]string{"publishable_key": app.Stripe.Key}, Data: map[string]any{"widget": widget}}, "stripe-js"); err != nil {
		app.ErrorLog.Println(err)
	}
}
