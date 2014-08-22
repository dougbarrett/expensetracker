package main

import (
	"database/sql"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/sessions"
	"log"
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

	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Map(SetupDB())

	m.Get("/", homepage)
	m.Get("/t/:hash", showTracker)
	m.Post("/t/save/:hash", saveTracker)
	m.Post("/i/new/:hash", saveItem)

	m.Run()
}
