package lib

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nayeemnishaat/go-web-app/api/model"
)

type Config struct {
	Port       int
	Env        string
	API        string
	RootRouter http.Handler
	Stripe     struct {
		Secret string
		Key    string
	}
	DB *model.SqlDB
}

type Application struct {
	Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
}
