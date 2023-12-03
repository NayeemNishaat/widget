package main

import (
	"encoding/json"
	"net/http"

	"github.com/nayeemnishaat/go-web-app/api/lib"
)

func (app *application) getPaymentIntent(w http.ResponseWriter, r *http.Request) {
	j := lib.Response{Error: false}
	out, err := json.MarshalIndent(j, "", "  ")

	if err != nil {
		app.ErrorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
