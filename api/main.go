package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/nayeemnishaat/go-web-app/api/lib"
	"github.com/nayeemnishaat/go-web-app/api/model"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	initApp()
	db := lib.InitDB()
	defer db.Close()
	app.DB = &model.SqlDB{Pool: db}

	err = serve(&app)
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
