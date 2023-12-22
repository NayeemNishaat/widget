package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nayeemnishaat/go-web-app/api/lib"
	"github.com/nayeemnishaat/go-web-app/api/model"
	webLib "github.com/nayeemnishaat/go-web-app/web/lib"
	"github.com/nayeemnishaat/go-web-app/web/template"
)

func (app *Application) TerminalPage(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "terminal", &template.TemplateData{StringMap: map[string]string{"publishable_key": app.Stripe.Key}}, "stripe-js"); err != nil {
		app.ErrorLog.Println(err)
	}
}

// GetTransactionData gets txn data from post and stripe
func (app *Application) GetTransactionData(r *http.Request) (webLib.TransactionData, error) {
	var txnData webLib.TransactionData

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
		return txnData, err
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	amount, _ := strconv.Atoi(paymentAmount)

	card := lib.Card{
		Secret: app.Config.Stripe.Secret,
		Key:    app.Config.Stripe.Key,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntent)
	if err != nil {
		app.ErrorLog.Println(err)
		return txnData, err
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.ErrorLog.Println(err)
		return txnData, err
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	txnData = webLib.TransactionData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		PaymentIntentID: paymentIntent,
		PaymentMethodID: paymentMethod,
		PaymentAmount:   amount,
		PaymentCurrency: paymentCurrency,
		LastFour:        lastFour,
		ExpiryMonth:     int(expiryMonth),
		ExpiryYear:      int(expiryYear),
		BankReturnCode:  pi.LatestCharge.ID,
	}
	return txnData, nil
}

func (app *Application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	widgetID, _ := strconv.Atoi(r.Form.Get("product_id"))

	txnData, err := app.GetTransactionData(r)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	// create a new customer
	customerID, err := app.SaveCustomer(txnData.FirstName, txnData.LastName, txnData.Email)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	// create a new transaction
	txn := model.Transactions{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
		TransactionStatusID: 7, // Note: 7 -> Cleared
	}

	txnID, err := app.SaveTransaction(txn)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	// create a new order
	order := model.Order{
		WidgetID:      widgetID,
		TransactionID: txnID,
		CustomerID:    customerID,
		StatusID:      4, // Note: 4 -> Cleared
		Quantity:      1,
		Amount:        txnData.PaymentAmount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, err = app.SaveOrder(order)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	app.Session.Put(r.Context(), "receipt", txnData)

	http.Redirect(w, r, "/ecom/receipt", http.StatusSeeOther)
}

// VirtualTerminalPaymentSucceeded displays the receipt page for virtual terminal transactions
func (app *Application) VirtualPaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	txnData, err := app.GetTransactionData(r)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	// create a new transaction
	txn := model.Transactions{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
		TransactionStatusID: 7,
	}

	_, err = app.SaveTransaction(txn)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	// write this data to session, and then redirect user to new page
	app.Session.Put(r.Context(), "receipt", txnData)
	http.Redirect(w, r, "/ecom/virtual-receipt", http.StatusSeeOther)
}

func (app *Application) Receipt(w http.ResponseWriter, r *http.Request) {
	txn := app.Session.Get(r.Context(), "receipt").(webLib.TransactionData)
	// app.Session.Remove(r.Context(), "receipt")

	if err := app.RenderTemplate(w, r, "receipt", &template.TemplateData{Data: map[string]any{"txn": txn}}); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) VirtualReceipt(w http.ResponseWriter, r *http.Request) {
	txn := app.Session.Get(r.Context(), "receipt").(webLib.TransactionData)
	// app.Session.Remove(r.Context(), "receipt")

	if err := app.RenderTemplate(w, r, "virtual-receipt", &template.TemplateData{Data: map[string]any{"txn": txn}}); err != nil {
		app.ErrorLog.Println(err)
	}
}

// SaveCustomer saves a customer and returns id
func (app *Application) SaveCustomer(firstName, lastName, email string) (int, error) {
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
func (app *Application) SaveTransaction(txn model.Transactions) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveOrder saves a order and returns id
func (app *Application) SaveOrder(order model.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

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

func (app *Application) BronzePlan(w http.ResponseWriter, r *http.Request) {
	widget, err := app.DB.GetWidget(2)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	data := make(map[string]any)
	data["widget"] = widget

	if err := app.RenderTemplate(w, r, "bronze-plan", &template.TemplateData{Data: data}); err != nil {
		app.ErrorLog.Println(err)
	}
}
