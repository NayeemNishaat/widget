package lib

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
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
	DB *pgxpool.Pool
}

type Application struct {
	Config
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
	Version       string
}
