package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"github.com/Flgado/bookings/pkg/config"
	"github.com/Flgado/bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData{
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// nget the template cache form the app config


	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	
	err := t.Execute(buf, td)

	if err != nil {
		log.Println(err)
	}
	//render the template
	_, err = buf.WriteTo(w)

	if err != nil {
		log.Println(err)
	}

	if err != nil {
		fmt.Println("error parsing template", err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myChache := map[string]*template.Template{}

	// get all of the files name *.page.html from ./templates
	pages, err := filepath.Glob("../../templates/*.page.html")
	if err != nil {
		return myChache, err
	}

	// tange through all files ending with *.page.html

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myChache, err
		}

		matches, err := filepath.Glob("../../templates/*layout.html")

		if err != nil {
			return myChache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("../../templates/*.layout.html")
		}

		myChache[name] = ts
	}

	return myChache, err

}
