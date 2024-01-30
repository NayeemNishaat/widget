package lib

var application *Application

func SetConfig(app *Application) {
	application = app
}

func GetConfig() *Application {
	return application
}
