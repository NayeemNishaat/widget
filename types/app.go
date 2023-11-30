package types

import (
	"html/template"
	"log"
	"net/http"
)

type Config struct {
	Port       int
	Env        string
	API        string
	RootRouter http.Handler
	DB         struct{ DSN string }
	Stripe     struct {
		Secret string
		Key    string
	}
}

type Application struct {
	Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
}
