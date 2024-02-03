package controller

import (
	"fmt"
	"net/http"

	"github.com/nayeemnishaat/go-web-app/api/lib"
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

func (app *Application) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "forgot-password", &template.TemplateData{}); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) ShowResetPassword(w http.ResponseWriter, r *http.Request) {
	theURL := r.RequestURI

	testURL := fmt.Sprintf("%s%s", app.FrontendURL, theURL)

	signer := lib.Signer{Secret: []byte(app.SigningSecret)}
	valid := signer.VerifyToken(testURL)

	if !valid {
		app.ErrorLog.Println("Invalid URL (Tampered)")
		return
	}

	// check expiry
	expired := signer.Expired(testURL, 60)
	if expired {
		app.ErrorLog.Println("Link Expired")
		return
	}

	data := make(map[string]any)
	data["email"] = r.URL.Query().Get("email")

	if err := app.RenderTemplate(w, r, "reset-password", &template.TemplateData{Data: data}); err != nil {
		app.ErrorLog.Println(err)
	}
}
