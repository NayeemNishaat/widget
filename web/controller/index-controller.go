package controller

import "github.com/nayeemnishaat/go-web-app/types"

type Application struct {
	*types.Application
}

var App *Application

func InitApp(app *types.Application) {
	App = &Application{app}
}
