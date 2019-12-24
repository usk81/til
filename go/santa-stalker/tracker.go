package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	statusTakeOff = "Takeoff"
	statusMove    = "Move"
	statusDeliver = "Deliver"
	statusFinsish = "Finsish"
)

type APIResponse struct {
	Status       string        `json:"status"`
	Language     string        `json:"language"`
	TimeOffset   int           `json:"timeOffset"`
	Fingerprint  string        `json:"fingerprint"`
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	ID                string   `json:"id"`
	Arrival           int64    `json:"arrival"`
	Departure         int64    `json:"departure"`
	Population        int      `json:"population"`
	PresentsDelivered int      `json:"presentsDelivered"`
	City              string   `json:"city"`
	Region            string   `json:"region"`
	Location          Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Tracker struct {
	Destinations []Destination
	Cache        Cache
}

type Cache struct {
	ExpiredAt   int64
	Status      string
	Destination Destination
}

type Result struct {
	IsCache     bool
	Status      string
	Destination Destination
}

func GetTracker(uri string) (result *Tracker, err error) {
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var r APIResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}

	return &Tracker{
		Destinations: r.Destinations,
	}, nil
}

func (tk *Tracker) SetCache(dest Destination, status string, expiredAt int64) {
	tk.Cache = Cache{
		Destination: dest,
		Status:      status,
		ExpiredAt:   expiredAt,
	}
}

func (tk *Tracker) Current(tt time.Time) Result {
	ts := tt.Unix() * 1000

	c := tk.Cache
	if c.Status != "" && c.ExpiredAt >= ts {
		return Result{
			Status:      c.Status,
			Destination: c.Destination,
			IsCache:     true,
		}
	}

	home := tk.Destinations[0]
	if home.Departure > ts {
		return Result{
			Status:      statusTakeOff,
			Destination: home,
		}
	}

	landing := tk.Destinations[len(tk.Destinations)-1]
	if landing.Arrival < ts {
		return Result{
			Status:      statusFinsish,
			Destination: landing,
		}
	}

	ds := tk.Destinations[1:]
	for _, dest := range ds {
		if dest.Arrival > ts {
			tk.SetCache(dest, statusMove, dest.Arrival-1)
			return Result{
				Status:      statusMove,
				Destination: dest,
			}
		} else if dest.Arrival <= ts && dest.Departure >= ts {
			tk.SetCache(home, statusDeliver, dest.Departure)
			return Result{
				Status:      statusDeliver,
				Destination: dest,
			}
		}
	}
	return Result{
		Status:      statusFinsish,
		Destination: landing,
	}
}
