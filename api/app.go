package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	Port       int
	Env        string
	RootRouter http.Handler
	DB         struct{ DSN string }
	Stripe     struct {
		Secret string
		Key    string
	}
}

type application struct {
	config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var app application

func initApp() {
	flag.IntVar(&app.Port, "port", 4000, "Server Port")
	flag.StringVar(&app.Env, "env", "dev", "App Env {dev|prod|maint}")

	flag.Parse()

	app.Stripe.Key = os.Getenv("STRIPE_KEY")
	app.Stripe.Secret = os.Getenv("STRIPE_SECRET")

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.RootRouter = Router()
}
