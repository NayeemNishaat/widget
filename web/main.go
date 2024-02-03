package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/nayeemnishaat/go-web-app/api/model"
	"github.com/nayeemnishaat/go-web-app/web/controller"
	"github.com/nayeemnishaat/go-web-app/web/lib"
	"github.com/nayeemnishaat/go-web-app/web/router"
	tmpl "github.com/nayeemnishaat/go-web-app/web/template"
)

const VERSION = "1.0.0"
const CSS_VERSION = "1"
const SESSION_LIFETIME = 24 * time.Hour

func main() {
	gob.Register(lib.TransactionData{}) // Important: Registering session variable type.

	var app lib.Application
	flag.IntVar(&app.Port, "port", 3000, "Server Port")
	flag.StringVar(&app.Env, "env", "dev", "App Env {dev|prod}")
	flag.StringVar(&app.API, "api", "http://localhost:4000", "API URL")
	flag.StringVar(&app.FrontendURL, "frontend_url", "http://localhost:3000", "Frontend Url")

	flag.Parse()

	app.Stripe.Key = os.Getenv("STRIPE_KEY")
	app.Stripe.Secret = os.Getenv("STRIPE_SECRET")

	app.SigningSecret = os.Getenv("SIGNING_SECRET")

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.TemplateCache = make(map[string]*template.Template)
	app.Version = VERSION

	db := lib.InitDB()
	defer db.Close()
	app.DB = &model.SqlDB{Pool: db}

	lib.Session = scs.New()
	lib.Session.Lifetime = SESSION_LIFETIME
	lib.Session.Store = pgxstore.New(db)
	app.Session = *lib.Session

	tApp := tmpl.InitApp(&app)
	controller.InitApp(tApp)

	app.RootRouter = router.RootRouter()

	lib.SetConfig(&app)

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
	if app.Env == "dev" {
		fmt.Printf("http://localhost:%d\n", app.Port)
	}

	return srv.ListenAndServe()
}
