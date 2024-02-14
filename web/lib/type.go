package lib

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/websocket"
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
	DB            *model.SqlDB
	Session       scs.SessionManager
	SigningSecret string
	FrontendURL   string
	MicroURL      string
}

type Application struct {
	Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
	WsChan        chan WsPayload
	Wg            *sync.WaitGroup
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

type Invoice struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type WsPayload struct {
	Action      string          `json:"action"`
	Message     string          `json:"message"`
	Username    string          `json:"username"`
	MessageType string          `json:"message_type"`
	UsedID      int             `json:"user_id"`
	Conn        *websocket.Conn `json:"-"`
}
