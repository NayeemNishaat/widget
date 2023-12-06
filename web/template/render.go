package template

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/nayeemnishaat/go-web-app/web/lib"
)

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]any
	CSRF            string
	Warning         string
	Info            string
	Error           string
	IsAuthenticated bool
	API             string
	CSSVersion      string
}

type Application struct {
	*lib.Application
}

var App *Application

func InitApp(app *lib.Application) *Application {
	return &Application{app}
}

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

func formatCurrency(n int) string {
	f := float32(n / 100)
	return fmt.Sprintf("$%.2f", f)
}

//go:embed *
var templateFS embed.FS

func (app *Application) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.API = app.API
	return td
}

func (app *Application) RenderTemplate(w http.ResponseWriter, r *http.Request, page string, td *TemplateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("%s.page.gohtml", page)

	_, templateInMap := app.TemplateCache[templateToRender]

	if app.Config.Env == "prod" && templateInMap {
		t = app.TemplateCache[templateToRender]
	} else {
		t, err = app.parseTemplate(partials, page, templateToRender)

		if err != nil {
			app.ErrorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &TemplateData{}
	}

	td = app.addDefaultData(td, r)

	err = t.Execute(w, td)

	if err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	return err
}

func (app *Application) parseTemplate(partials []string, page string, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("%s.partial.gohtml", x)
		}
	}

	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "base.layout.gohtml", strings.Join(partials, ","), templateToRender)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "base.layout.gohtml", templateToRender)
	}

	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}

	app.TemplateCache[templateToRender] = t

	return t, nil
}
