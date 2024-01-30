package controller

import (
	"net/http"

	"github.com/nayeemnishaat/go-web-app/web/template"
)

func (app *Application) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "login", &template.TemplateData{}); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	app.Session.RenewToken(r.Context())
	if err := r.ParseForm(); err != nil {
		app.ErrorLog.Println(err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := app.DB.Authenticate(email, password)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
