package main

import (
	"database/sql"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"net/http"
	"strconv"
	"time"
)

func homepage(db *sql.DB, w http.ResponseWriter, r *http.Request) string {

	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Header().Set("Expires", "0")
	var tracker Tracker
	tracker.Setup(db)
	tracker.New()

	if tracker.Hash == "" {
		return "We're having a bit of an issue taking you to where you want to go, try again later"
	}

	http.Redirect(w, r, "/t/"+tracker.Hash, http.StatusFound)
	return "Uh oh! Something went wrong, please try again later..."
}

func showTracker(db *sql.DB, params martini.Params, r render.Render, session sessions.Session) {
	var retData struct {
		Tracker
		SavedItem bool
	}
	retData.Setup(db)
	retData.Hash = params["hash"]
	retData.FindByHash()

	if retData.ID != "" {
		v := session.Get("savedItem")
		if v == nil {
			retData.SavedItem = false
		} else {
			retData.SavedItem = true
			session.Delete("savedItem")
		}
		r.HTML(200, "tracker", retData)
	} else {
		r.HTML(404, "tracker_notfound", nil)
	}
}

func saveTracker(db *sql.DB, params martini.Params, r render.Render, req *http.Request) {
	var tracker Tracker
	tracker.Setup(db)
	tracker.Hash = params["hash"]
	tracker.FindByHash()

	if tracker.ID != req.FormValue("id") {
		r.HTML(500, "tracker_saveerror", nil)
		return
	}
	tracker.Name = req.FormValue("name")
	tracker.Save()

	r.Redirect("/t/" + tracker.Hash)
}

func saveItem(db *sql.DB, params martini.Params, r render.Render, req *http.Request, session sessions.Session) {
	var tracker Tracker
	tracker.Setup(db)
	tracker.Hash = params["hash"]
	tracker.FindByHash()

	if tracker.ID != req.FormValue("id") {
		r.HTML(500, "tracker_saveerror", nil)
		return
	}

	var item Item
	item.Label = req.FormValue("label")
	if req.FormValue("spend") != "" {
		item.Spend, _ = strconv.ParseFloat(req.FormValue("spend"), 64)
	}
	if req.FormValue("save") != "" {
		item.Save, _ = strconv.ParseFloat(req.FormValue("save"), 64)
	}
	item.TrackerID = tracker.ID
	item.Timeslice = time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		0, 0, 0, 0, time.UTC)

	if err := tracker.SaveItem(&item); err != nil {
		r.HTML(500, "tracker_saveerror", err.Error())
	}

	session.Set("savedItem", "true")
	r.Redirect("/t/" + tracker.Hash)
}
