package controller

import (
	"net/http"

	"github.com/nayeemnishaat/go-web-app/web/template"
)

func (app *Application) TerminalPage(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "terminal", &template.TemplateData{StringMap: map[string]string{"publishable_key": app.Stripe.Key}}, "stripe-js"); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) AllSales(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "all-sales", &template.TemplateData{}); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) AllSubscriptions(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "all-subscriptions", &template.TemplateData{}); err != nil {
		app.ErrorLog.Println(err)
	}
}
