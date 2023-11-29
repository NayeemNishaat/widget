package main

import (
	"flag"
	"html/template"
	"log"
	"os"

	"github.com/nayeemnishaat/go-web-app/controller"
	"github.com/nayeemnishaat/go-web-app/lib"
	"github.com/nayeemnishaat/go-web-app/router"
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

	lib.InitApp(&app)
	controller.InitApp(&app)

	app.RootRouter = router.RootRouter()

	err := lib.App.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
