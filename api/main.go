package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nayeemnishaat/go-web-app/api/lib"
)

func main() {
	initApp()
	db := lib.InitDB()
	defer db.Close()

	err := serve(&app)
	if err != nil {
		log.Fatal(err)
	}
}

func serve(app *application) error {
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
