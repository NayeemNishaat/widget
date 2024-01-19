package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/api/lib"
)

func (app *application) getPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload lib.StripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	card := lib.Card{Secret: app.Stripe.Secret, Key: app.Stripe.Key, Currency: payload.Currency}

	okay := true
	pi, msg, err := card.CreatePaymentIntent(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(pi, "", "  ")
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := lib.Response{Error: true, Message: msg, Data: map[string]any{}}
		out, err := json.MarshalIndent(j, "", "  ")

		if err != nil {
			app.ErrorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

func (app *application) getWidgetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(widget, "", "  ") // Warning: Don't use indent in production.
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) createCustomerAndSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	var data lib.StripePayload
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	okay := true
	msg := ""

	out, err := json.MarshalIndent(map[string]any{"OK": okay, "Message": msg}, "", "  ")

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
