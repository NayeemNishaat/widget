package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/api/lib"
	"github.com/nayeemnishaat/go-web-app/api/model"
	"github.com/stripe/stripe-go/v76"
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

// SaveCustomer saves a customer and returns id
func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := model.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveTransaction saves a txn and returns id
func (app *application) SaveTxn(txn model.Transactions) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveOrder saves a order and returns id
func (app *application) SaveOrder(order model.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app *application) createCustomerAndSubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	var data lib.StripePayload
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	card := lib.Card{Secret: app.Stripe.Secret, Key: app.Stripe.Key, Currency: data.Currency}

	okay := true
	var subscription *stripe.Subscription
	txnMsg := "Transaction Successful."

	stripeCustomer, msg, err := card.CreateCustomer(data.PaymentMethod, data.Email)
	if err != nil {
		app.ErrorLog.Println(err)
		okay = false
		txnMsg = msg
	}

	if okay {
		subscription, err = card.SubscribeToPlan(stripeCustomer, data.Plan, data.Email, data.LastFour, "")
		if err != nil {
			app.ErrorLog.Println(err)
			okay = false
			txnMsg = "Error subscribing customer"
		}

		app.InfoLog.Println(subscription.ID)
	}

	if okay {
		productID, _ := strconv.Atoi(data.ProductID)
		customerID, err := app.SaveCustomer(data.FirstName, data.LastName, data.Email)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		// Create a new txn
		amount, _ := strconv.Atoi(data.Amount)
		// expiryMonth, _ := strconv.Atoi(data.ExpiryMonth)
		// expiryYear, _ := strconv.Atoi(data.ExpiryYear)
		txn := model.Transactions{
			Amount:              amount,
			Currency:            "usd",
			LastFour:            data.LastFour,
			ExpiryMonth:         data.ExpiryMonth,
			ExpiryYear:          data.ExpiryYear,
			TransactionStatusID: 7,
		}

		txnID, err := app.SaveTxn(txn)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		// Create Order
		order := model.Order{WidgetID: productID, TransactionID: txnID, CustomerID: customerID, StatusID: 4, Quantity: 1, Amount: amount, CreatedAt: time.Now(), UpdatedAt: time.Now()}

		_, err = app.SaveOrder(order)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}
	}

	out, err := json.MarshalIndent(map[string]any{"OK": okay, "Message": txnMsg}, "", "  ")

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
