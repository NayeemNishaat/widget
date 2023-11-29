package controller

import (
	"fmt"
	"net/http"
)

func (app *Application) User(w http.ResponseWriter, r *http.Request) {
	fmt.Println("OK")
}
