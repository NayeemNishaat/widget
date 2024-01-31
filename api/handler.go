package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

func (app *application) createAuthToken(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := lib.ReadJSON(w, r, &userInput)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		lib.InvalidCredentials(w)
		return
	}

	// Point: Validate Password
	validPassword, err := lib.PasswordMatches(user.Password, userInput.Password)
	if err != nil {
		lib.InvalidCredentials(w)
		return
	}

	if !validPassword {
		lib.InvalidCredentials(w)
		return
	}

	// Point: Generate Token
	token, err := model.GenerateToken(user.ID, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	// Point: Save to DB
	err = app.DB.InsertToken(token, user)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	var payload struct {
		Error   bool         `json:"error"`
		Message string       `json:"mesage"`
		Token   *model.Token `json:"authentication_token"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("Token for %s is created.", userInput.Email)
	payload.Token = token

	// out, err := json.MarshalIndent(payload, "", "\t")
	// if err != nil {
	// 	app.ErrorLog.Println(err)
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(out)

	_ = lib.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) authenticateToken(r *http.Request) (*model.User, error) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authorization token wrong size")
	}

	user, err := app.DB.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

func (app *application) checkAuthentication(w http.ResponseWriter, r *http.Request) {
	user, err := app.authenticateToken(r)
	if err != nil {
		lib.InvalidCredentials(w)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{Error: false, Message: fmt.Sprintf("authenticated user %s", user.Email)}
	lib.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) VirtualTerminalPaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	var txnData struct {
		PaymentAmount   int    `json:"amount"`
		PaymentCurrency string `json:"currency"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		PaymentIntent   string `json:"payment_intent"`
		PaymentMethod   string `json:"payment_method"`
		BankReturnCode  string `json:"bank_return_code"`
		ExpiryMonth     int    `json:"expiry_month"`
		ExpiryYear      int    `json:"expiry_year"`
		LastFour        string `json:"last_four"`
	}

	err := lib.ReadJSON(w, r, &txnData)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	card := lib.Card{
		Secret: app.Stripe.Secret,
		Key:    app.Stripe.Secret,
	}

	pi, err := card.RetrievePaymentIntent(txnData.PaymentIntent)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	pm, err := card.GetPaymentMethod(txnData.PaymentMethod)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	txnData.LastFour = pm.Card.Last4
	txnData.ExpiryMonth = int(pm.Card.ExpMonth)
	txnData.ExpiryYear = int(pm.Card.ExpYear)

	txn := model.Transactions{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      pi.LatestCharge.ID,
		PaymentIntent:       txnData.PaymentIntent,
		PaymentMethod:       txnData.PaymentMethod,
		TransactionStatusID: 7,
	}

	_, err = app.SaveTxn(txn)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	lib.WriteJSON(w, http.StatusOK, txn)
}

func (app *application) sendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}

	err := lib.ReadJSON(w, r, &payload)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	var data struct {
		Link string
	}
	data.Link = "http://www.unb.ca"

	// send mail
}
