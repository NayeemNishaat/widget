package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	Port int
	SMTP struct {
		Host     string
		Username string
		Password string
		Port     int
	}
	FrontendURL string
}

type application struct {
	config
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Version  string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.Port, "port", 5000, "Invoice Microservice Port")
	flag.StringVar(&cfg.SMTP.Host, "smtp_host", "sandbox.smtp.mailtrap.io", "SMTP Host")
	flag.StringVar(&cfg.SMTP.Username, "smtp_username", "1d69f754aa40ea", "SMTP Username")
	flag.StringVar(&cfg.SMTP.Password, "smtp_password", "31b59553edd95d", "SMTP Password")
	flag.IntVar(&cfg.SMTP.Port, "smtp_port", 25, "SMTP Port")
	flag.StringVar(&cfg.FrontendURL, "frontend_url", "http://localhost:3000", "Frontend Url")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   cfg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Version:  version,
	}

	CreateDirIfNotExist("./microservice/invoice/pdf")

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	srv := http.Server{
		Addr:              fmt.Sprintf(":%d", app.Port),
		Handler:           app.Router(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.InfoLog.Printf("Invoice microservice running on port %d", app.Port)

	return srv.ListenAndServe()
}
