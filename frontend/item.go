package main

import (
	"time"
)

type Item struct {
	Id        string
	Label     string
	TrackerID string
	Spend     float64
	Save      float64
	Timeslice time.Time
}
