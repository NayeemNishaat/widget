package controller

import (
	"github.com/nayeemnishaat/go-web-app/web/template"
)

type Application struct {
	*template.Application
}

var App *Application

func InitApp(app *template.Application) {
	App = &Application{app}
}
