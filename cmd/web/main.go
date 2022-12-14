package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Flgado/bookings/pkg/config"
	"github.com/Flgado/bookings/pkg/handlers"
	"github.com/Flgado/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
//main is the main application function
func main() {
	

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc , err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	// start the repository
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// give to render the access to App configuration 
	render.NewTemplates(&app)


	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	//http.ListenAndServe(portNumber, nil)

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
