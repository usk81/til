package main

import (
	"fmt"
	"math"
)

const radius = 6378137.0

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func main() {
	from := Location{
		Lat: 35.693825,
		Lng: 139.703356,
	}
	to := Location{
		Lat: 35.644816899999995,
		Lng: 139.6980594,
	}

	fmt.Println(distance(from, to))
}

func distance(from, to Location) float64 {
	fx := deg2rad(from.Lat)
	fy := deg2rad(from.Lng)

	tx := deg2rad(to.Lat)
	ty := deg2rad(to.Lng)

	averageLat := (fx - tx) / 2
	averageLon := (fy - ty) / 2

	return radius * 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(averageLat), 2)+(math.Cos(fx)*math.Cos(tx)*math.Pow(math.Sin(averageLon), 2))))
}

// deg2rad transforms radical value
func deg2rad(r float64) float64 {
	return (r * math.Pi) / 180.0
}
