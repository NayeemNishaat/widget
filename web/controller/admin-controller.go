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

func (app *Application) ShowSale(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Sale"
	stringMap["cancel"] = "/admin/all-sales"
	stringMap["refund-url"] = "/api/v1/admin/refund"
	stringMap["refund-btn"] = "Refund Order"
	stringMap["refunded-badge"] = "Refunded"
	stringMap["refunded-msg"] = "Charge Refunded"

	if err := app.RenderTemplate(w, r, "sale", &template.TemplateData{StringMap: stringMap}); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) ShowSubscription(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["title"] = "Subscription"
	stringMap["cancel"] = "/admin/all-subscriptions"
	stringMap["refund-url"] = "/api/v1/admin/cancel-subscription"
	stringMap["refund-btn"] = "Cancel Subscription"
	stringMap["refunded-badge"] = "Cancelled"
	stringMap["refunded-msg"] = "Subscription Cancelled"

	if err := app.RenderTemplate(w, r, "sale", &template.TemplateData{StringMap: stringMap}); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) AllUsers(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "all-users", &template.TemplateData{}); err != nil {
		app.ErrorLog.Println(err)
	}
}

// Sows one user's details
func (app *Application) OneUser(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "one-user", &template.TemplateData{}); err != nil {
		app.ErrorLog.Println(err)
	}
}
