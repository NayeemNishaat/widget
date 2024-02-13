package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	// "github.com/nayeemnishaat/go-web-app/api/lib"
	// "github.com/nayeemnishaat/go-web-app/api/model"
)

var wg sync.WaitGroup

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	initApp()
	// db := lib.InitDB()
	// defer db.Close()
	// app.DB = &model.SqlDB{Pool: db}
	app.Wg = &wg

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
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	app.InfoLog.Printf("Starting server on %s mode on port %d\n", app.Env, app.Port)

	go gracefulShutdown()
	return srv.ListenAndServe()
}

var quitChan = make(chan os.Signal, 1) // Important: Must be a buffered chan because we are waiting for quit below. If it's not buffered then signal will not be sent to the chan because listener is not listening at the same time. So for unbuffered chan both sender and receiver should be ready at the same time via separate goroutines.
func gracefulShutdown() {
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	<-quitChan

	shutdown()
	os.Exit(0)
}

func shutdown() {
	fmt.Println("\nPerforming Cleanup")

	wg.Wait()
	close(quitChan)
}
