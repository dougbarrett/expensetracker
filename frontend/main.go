package main

import (
	"database/sql"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
)

var config tomlConfig

func SetupDB() *sql.DB {
	queryUrl := config.Mysql.Username + ":" + config.Mysql.Password + "@tcp(" + config.Mysql.Host + ":" + config.Mysql.Port + ")/" + config.Mysql.Database
	db, err := sql.Open("mysql", queryUrl)

	if err != nil {
		panic(err)
	}

	return db
}

func main() {

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
		return
	}
	log.Printf("Session Name: %s", config.Session.Name)
	log.Printf("Session Key: %s", config.Session.Key)
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte(config.Session.Key))
	m.Use(sessions.Sessions(config.Session.Name, store))

	m.Use(csrf.Generate(&csrf.Options{
		Secret:     config.Session.Key + config.Session.Name,
		SessionKey: "trackingHash",
		// Custom error response.
		ErrorFunc: func(w http.ResponseWriter) {
			http.Error(w, "CSRF token validation failed", http.StatusBadRequest)
		},
	}))

	m.Use(gzip.All())

	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Map(SetupDB())

	m.Get("/", homepage)
	m.Get("/t/:hash", showTracker)
	m.Post("/t/save/:hash", csrf.Validate, saveTracker)
	m.Post("/i/new/:hash", csrf.Validate, saveItem)

	m.Run()
}
