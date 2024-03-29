package main

import (
	"bytes"
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
	"golang.org/x/crypto/bcrypt"
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
func (app *application) saveCustomer(firstName, lastName, email string) (int, error) {
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
func (app *application) saveTxn(txn model.Transactions) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveOrder saves a order and returns id
func (app *application) saveOrder(order model.Order) (int, error) {
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

	// validate data
	v := lib.NewValidator()
	v.Check(len(data.FirstName) > 1, "first_name", "must be at least 2 characters")
	v.Check(data.LastName != "", "last_name", "cannot be empty")

	if !v.Valid() {
		app.fieldValidation(w, r, v.Errors)
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
	}

	if okay {
		productID, _ := strconv.Atoi(data.ProductID)
		customerID, err := app.saveCustomer(data.FirstName, data.LastName, data.Email)
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
			PaymentIntent:       subscription.ID,
			PaymentMethod:       data.PaymentMethod,
		}

		txnID, err := app.saveTxn(txn)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		// Create Order
		order := model.Order{WidgetID: productID, TransactionID: txnID, CustomerID: customerID, StatusID: 4, Quantity: 1, Amount: amount, CreatedAt: time.Now(), UpdatedAt: time.Now()}

		orderID, err := app.saveOrder(order)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		inv := lib.Invoice{
			ID:        orderID,
			Amount:    2000,
			Product:   "Bronze Plan monthly subscription",
			Quantity:  order.Quantity,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
			CreatedAt: time.Now(),
		}

		app.Wg.Add(1)
		go app.callInvoiceMicro(inv)
		// err = app.callInvoiceMicro(inv)
		// if err != nil {
		// 	app.ErrorLog.Println(err)
		// }
	}

	out, err := json.MarshalIndent(map[string]any{"error": !okay, "message": txnMsg}, "", "  ")

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// callInvoiceMicro calls the invoicing microservice
func (app *application) callInvoiceMicro(inv lib.Invoice) error {
	defer app.Wg.Done()

	url := "http://localhost:5000/invoice/create-and-send"
	out, err := json.MarshalIndent(inv, "", "\t")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
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

func (app *application) virtualTerminalPaymentSucceeded(w http.ResponseWriter, r *http.Request) {
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

	_, err = app.saveTxn(txn)
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

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	// verify email exists
	_, err = app.DB.GetUserByEmail(payload.Email)
	if err != nil {
		resp.Error = true
		resp.Message = "Email doesn't exist!"

		lib.WriteJSON(w, http.StatusAccepted, resp)
		return
	}

	link := fmt.Sprintf("%s/auth/reset-password?email=%s", app.config.FrontendURL, payload.Email)
	sign := lib.Signer{Secret: []byte(app.config.SigningSecret)}

	signedLink := sign.GenerateTokenFromString(link)

	var data struct {
		Link string
	}
	data.Link = signedLink

	// send mail
	err = app.SendMail("info@widgets.com", payload.Email, "Password Reset Request", "password-reset", data)
	if err != nil {
		app.ErrorLog.Println(err)
		lib.BadRequest(w, r, err)
		return
	}

	resp.Error = false
	resp.Message = "Success"

	lib.WriteJSON(w, http.StatusCreated, resp)
}

func (app *application) resetPassword(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := lib.ReadJSON(w, r, &payload)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	encryptor := lib.Encryption{
		Key: []byte(app.SigningSecret),
	}

	realEmail, err := encryptor.Decrypt(payload.Email)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	user, err := app.DB.GetUserByEmail(realEmail)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	err = app.DB.UpdatePassword(user, string(newHash))
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "Password Updated!"

	lib.WriteJSON(w, http.StatusCreated, resp)
}

func (app *application) allSales(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"page"`
	}

	err := lib.ReadJSON(w, r, &payload)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	allSales, lastPage, totalRecords, err := app.DB.GetAllOrdersPaginated(payload.PageSize, payload.CurrentPage)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	var resp struct {
		CurrentPage  int            `json:"current_page"`
		PageSize     int            `json:"page_size"`
		LastPage     int            `json:"last_page"`
		TotalRecords int            `json:"total_records"`
		Orders       []*model.Order `json:"orders"`
	}

	resp.CurrentPage = payload.CurrentPage
	resp.PageSize = payload.PageSize
	resp.LastPage = lastPage
	resp.TotalRecords = totalRecords
	resp.Orders = allSales

	lib.WriteJSON(w, http.StatusOK, resp)
}

// AllSubscriptions returns all subscriptions as a slice
func (app *application) allSubscriptions(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PageSize    int `json:"page_size"`
		CurrentPage int `json:"page"`
	}

	err := lib.ReadJSON(w, r, &payload)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	allSales, lastPage, totalRecords, err := app.DB.GetAllSubscriptionsPaginated(payload.PageSize, payload.CurrentPage)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	var resp struct {
		CurrentPage  int            `json:"current_page"`
		PageSize     int            `json:"page_size"`
		LastPage     int            `json:"last_page"`
		TotalRecords int            `json:"total_records"`
		Orders       []*model.Order `json:"orders"`
	}

	resp.CurrentPage = payload.CurrentPage
	resp.PageSize = payload.PageSize
	resp.LastPage = lastPage
	resp.TotalRecords = totalRecords
	resp.Orders = allSales

	lib.WriteJSON(w, http.StatusOK, resp)
}

func (app *application) getSale(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orderID, _ := strconv.Atoi(id)

	order, err := app.DB.GetOrderByID(orderID)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	lib.WriteJSON(w, http.StatusOK, order)
}

func (app *application) refundCharge(w http.ResponseWriter, r *http.Request) {
	var chargeToRefund struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"`
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
	}

	err := lib.ReadJSON(w, r, &chargeToRefund)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	// validate amount

	card := lib.Card{
		Secret:   app.Stripe.Secret,
		Key:      app.Stripe.Key,
		Currency: chargeToRefund.Currency,
	}

	err = card.Refund(chargeToRefund.PaymentIntent, chargeToRefund.Amount)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	// update status in db
	err = app.DB.UpdateOrderStatus(chargeToRefund.ID, 5)
	if err != nil {
		lib.BadRequest(w, r, errors.New("the charge was refunded but the db could not be updated"))
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "Charge Refunded"

	lib.WriteJSON(w, http.StatusCreated, resp)
}

func (app *application) cancelSubscription(w http.ResponseWriter, r *http.Request) {
	var subToCancel struct {
		ID            int    `json:"id"`
		PaymentIntent string `json:"pi"`
		Currency      string `json:"currency"`
	}

	err := lib.ReadJSON(w, r, &subToCancel)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	card := lib.Card{
		Secret:   app.config.Stripe.Secret,
		Key:      app.config.Stripe.Key,
		Currency: subToCancel.Currency,
	}

	err = card.CancelSubscription(subToCancel.PaymentIntent)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	// update status in db
	err = app.DB.UpdateOrderStatus(subToCancel.ID, 6)
	if err != nil {
		lib.BadRequest(w, r, errors.New("the subscription was cancelled but the db could not be updated"))
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = "Subscription Cancelled"

	lib.WriteJSON(w, http.StatusCreated, resp)
}

// AllUsers returns a JSON file listing all admin users
func (app *application) allUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := app.DB.GetAllUsers()
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	lib.WriteJSON(w, http.StatusOK, allUsers)
}

// OneUser gets one user by id (from the url) and returns it as JSON
func (app *application) oneUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	user, err := app.DB.GetOneUser(userID)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	lib.WriteJSON(w, http.StatusOK, user)
}

// EditUser is the handler for adding or editing an existing user
func (app *application) editUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	var user model.User

	err := lib.ReadJSON(w, r, &user)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	if userID > 0 {
		err = app.DB.EditUser(user)
		if err != nil {
			lib.BadRequest(w, r, err)
			return
		}

		if user.Password != "" {
			newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			if err != nil {
				lib.BadRequest(w, r, err)
				return
			}

			err = app.DB.UpdatePassword(user, string(newHash))
			if err != nil {
				lib.BadRequest(w, r, err)
				return
			}
		}
	} else {
		newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			lib.BadRequest(w, r, err)
			return
		}
		err = app.DB.AddUser(user, string(newHash))
		if err != nil {
			lib.BadRequest(w, r, err)
			return
		}
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	lib.WriteJSON(w, http.StatusOK, resp)
}

// DeleteUser deletes a user, and all associated tokens, from the database
func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	err := app.DB.DeleteUser(userID)
	if err != nil {
		lib.BadRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	lib.WriteJSON(w, http.StatusOK, resp)
}

func (app *application) fieldValidation(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	var payload struct {
		Error   bool              `json:"error"`
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors"`
	}

	payload.Error = true
	payload.Message = "failed validation"
	payload.Errors = errors
	lib.WriteJSON(w, http.StatusUnprocessableEntity, payload)
}
