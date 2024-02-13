package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/nayeemnishaat/go-web-app/api/model"
)

type config struct {
	Port       int
	Env        string
	RootRouter http.Handler
	Stripe     struct {
		Secret string
		Key    string
	}
	DB   *model.SqlDB
	Wg   *sync.WaitGroup
	SMTP struct {
		Host     string
		Username string
		Password string
		Port     int
	}
	SigningSecret string
	FrontendURL   string
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
	flag.StringVar(&app.FrontendURL, "frontend_url", "http://localhost:3000", "Frontend Url")

	flag.Parse()

	app.Stripe.Key = os.Getenv("STRIPE_KEY")
	app.Stripe.Secret = os.Getenv("STRIPE_SECRET")

	app.SMTP.Host = os.Getenv("SMTP_HOST")
	app.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	intPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Panic(err)
	} else {
		app.SMTP.Port = intPort
	}
	app.SMTP.Username = os.Getenv("SMTP_USERNAME")

	app.SigningSecret = os.Getenv("SIGNING_SECRET")

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.RootRouter = Router()
}
