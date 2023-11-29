package controller

import "github.com/nayeemnishaat/go-web-app/lib"

type Application struct {
	*lib.Application
}

var App *Application

func InitApp(app *lib.Application) {
	App = &Application{app}
}
