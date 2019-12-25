package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ikasamah/homecast"
)

// km/h
const speed = 44.6

func main() {
	l, err := env2Location()
	if err != nil {
		panic(err)
	}
	ss, err := getStationsFromCSV("station.csv")
	if err != nil {
		panic(err)
	}

	_, dist, err := nearestStation(l, ss)
	if err != nil {
		panic(err)
	}

	ms := speed * 1000.0 / (60.0 * 60.0)

	googlehome("私、メリーさん。今、近くの駅にいるの", "ja")
	time.Sleep(5 * time.Second)
	var flg bool
	for {
		switch {
		case dist <= 0:
			googlehome("私、メリーさん。今、あなたのお家の前にいるの", "ja")
		case dist <= 50 && !flg:
			googlehome("私、メリーさん。今、あなたのお家の近くにいるの。", "ja")
			flg = true
		}
		if dist <= 0 {
			break
		}
		time.Sleep(1 * time.Second)
		dist -= ms
		fmt.Println(dist)
	}
	googlehome("私、メリーさん。今、あなたの後ろにいるの。", "ja")
	googlehome("Hahaha, It's joke!!", "en")
}

func env2Location() (loc Location, err error) {
	strLoc := os.Getenv("Location")
	if strLoc == "" {
		err = errors.New("Can't get location")
		return
	}
	ls := strings.Split(strLoc, ",")
	if len(ls) != 2 {
		err = errors.New(strLoc + " is invalid location")
	}
	lat, err := strconv.ParseFloat(ls[0], 64)
	if err != nil {
		return
	}
	lng, err := strconv.ParseFloat(ls[1], 64)
	if err != nil {
		return
	}
	return Location{
		Lat: lat,
		Lng: lng,
	}, nil
}

func googlehome(msg, lang string) {
	fmt.Println(msg)
	ctx := context.Background()
	devices := homecast.LookupAndConnect(ctx)
	for _, device := range devices {
		device.Speak(ctx, msg, lang)
	}
	time.Sleep(5 * time.Second)
}
