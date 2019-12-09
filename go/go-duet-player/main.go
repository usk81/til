// package main plays two audio file
//   failure: oto.NewContext can be called only once
package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"

	"github.com/usk81/til/go-duet-player/speaker"
)

func main() {
	f1, err := os.Open("test1.mp3")
	if err != nil {
		log.Fatal(err)
	}
	st1, format, err := mp3.Decode(f1)
	if err != nil {
		log.Fatal(err)
	}
	defer st1.Close()

	sp1, err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	done1 := make(chan bool)

	f2, err := os.Open("test2.mp3")
	if err != nil {
		log.Fatal(err)
	}
	st2, format, err := mp3.Decode(f2)
	if err != nil {
		log.Fatal(err)
	}
	defer st2.Close()

	sp2, err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	done2 := make(chan bool)

	sp1.Play(beep.Seq(st1, beep.Callback(func() {
		done1 <- true
	})))
	sp2.Play(beep.Seq(st2, beep.Callback(func() {
		done2 <- true
	})))
	<-done1
	<-done2
}
