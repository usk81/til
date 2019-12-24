package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ikasamah/homecast"
)

const (
	myID     = "tokyo"
	santaURL = "https://firebasestorage.googleapis.com/v0/b/santa-tracker-firebase.appspot.com/o/route%2Fsanta_en.json?alt=media&2019b"
)

var jst *time.Location

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
}

func main() {
	tk, err := GetTracker(santaURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Start stalking")
	for {
		r := tk.Current(time.Now())
		if !r.IsCache {
			if r.Destination.ID == myID {
				googlehome("Santa Claus is coming to town!! Santa Claus is coming to town!! Santa Claus is coming to town!!")
				break
			} else {
				if r.Status == statusDeliver {
					googlehome(fmt.Sprintf("Santa Claus is in %s, %s\n", r.Destination.City, r.Destination.Region))
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("finish")
}

func googlehome(msg string) {
	fmt.Println(msg)
	ctx := context.Background()
	devices := homecast.LookupAndConnect(ctx)
	for _, device := range devices {
		device.Speak(ctx, msg, "en")
	}
}
