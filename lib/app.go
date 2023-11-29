package lib

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
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

var App *Application

func InitApp(app *Application) {
	App = app
}

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.Port),
		Handler:           app.RootRouter,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.InfoLog.Printf("Starting server on %s mode on port %d\n", app.Env, app.Port)

	return srv.ListenAndServe()
}
