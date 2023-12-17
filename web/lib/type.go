package lib

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/nayeemnishaat/go-web-app/api/model"
)

var Session *scs.SessionManager

type Config struct {
	Port       int
	Env        string
	API        string
	RootRouter http.Handler
	Stripe     struct {
		Secret string
		Key    string
	}
	DB      *model.SqlDB
	Session scs.SessionManager
}

type Application struct {
	Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
}

type TransactionData struct {
	FirstName       string
	LastName        string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	PaymentAmount   int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}
