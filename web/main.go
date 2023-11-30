package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nayeemnishaat/go-web-app/lib"
	"github.com/nayeemnishaat/go-web-app/web/controller"
	"github.com/nayeemnishaat/go-web-app/web/router"
)

const VERSION = "1.0.0"
const CSS_VERSION = "1"

func main() {
	var app lib.Application
	flag.IntVar(&app.Port, "port", 3000, "Server Port")
	flag.StringVar(&app.Env, "env", "dev", "App Env {dev|prod}")
	flag.StringVar(&app.API, "api", "http://localhost:4000", "API URL")

	flag.Parse()

	app.Stripe.Key = os.Getenv("STRIPE_KEY")
	app.Stripe.Secret = os.Getenv("STRIPE_SECRET")

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.TemplateCache = make(map[string]*template.Template)
	app.Version = VERSION

	controller.InitApp(&app)

	app.RootRouter = router.RootRouter()

	err := Serve(&app)
	if err != nil {
		log.Fatal(err)
	}
}

func Serve(app *lib.Application) error {
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
