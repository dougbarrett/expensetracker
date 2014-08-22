package main

import (
	"database/sql"
	"fmt"
	"github.com/dougbarrett/expensetracker/helpers"
	"time"
)

type Tracker struct {
	ID    string
	Name  string
	Hash  string
	db    *sql.DB
	Day   float64
	Month float64
	All   float64
}

func (t *Tracker) Setup(db *sql.DB) {
	t.db = db
}

func (t *Tracker) New() {

	for t.ID == "" && t.Hash == "" {
		hash := helpers.RandString(16)
		t.Hash = hash
		t.FindByHash()
		fmt.Println("searched for %s", t.Hash)
	}

	t.Create()
}

func (t *Tracker) FindByHash() {
	var today = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	var month = time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC)
	t.db.QueryRow("select id, name, hash from tracker where hash = ?", t.Hash).Scan(&t.ID, &t.Name, &t.Hash)
	t.db.QueryRow("select sum(save) - sum(spend) from item where tracker_id = ? and timeslice = ?", t.ID, today).Scan(&t.Day)
	t.db.QueryRow("select sum(save) - sum(spend) from item where tracker_id = ? and timeslice >= ? and timeslice <= ?", t.ID, month, today).Scan(&t.Month)
	t.db.QueryRow("select sum(save) - sum(spend) from item where tracker_id = ?", t.ID).Scan(&t.All)
}

func (t *Tracker) Create() {
	_, err := t.db.Exec("insert into tracker (`hash`) VALUES(?)", t.Hash)

	if err != nil {
		fmt.Println("err", err.Error())
	}
}

func (t *Tracker) Save() {
	fmt.Println("name", t.Name)
	fmt.Println("hash", t.Hash)
	_, err := t.db.Exec("update tracker SET `name` = ? WHERE hash = ?", t.Name, t.Hash)

	if err != nil {
		fmt.Println("err", err.Error())
	}
}

func (t *Tracker) SaveItem(item *Item) error {
	fmt.Println("item", item)
	_, err := t.db.Exec("insert into item (`tracker_id`, `label`, `spend`, `save`, `timeslice`) VALUES (?, ?, ?, ?, ?)", t.ID, item.Label, item.Spend, item.Save, item.Timeslice)

	return err
}
