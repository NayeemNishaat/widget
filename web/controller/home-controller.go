package controller

import (
	"net/http"

	"github.com/nayeemnishaat/go-web-app/web/template"
)

func (app *Application) HomePage(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "home", &template.TemplateData{}); err != nil {
		app.ErrorLog.Println(err)
	}
}
