package speaker

import (
	"log"
	"sync"

	"github.com/faiface/beep"
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"
)

type Player struct {
	mu      sync.Mutex
	mixer   beep.Mixer
	samples [][2]float64
	buf     []byte
	context *oto.Context
	player  *oto.Player
	done    chan struct{}
}

// Init initializes audio playback through speaker. Must be called before using this package.
//
// The bufferSize argument specifies the number of samples of the speaker's buffer. Bigger
// bufferSize means lower CPU usage and more reliable playback. Lower bufferSize means better
// responsiveness and less delay.
func Init(sampleRate beep.SampleRate, bufferSize int) (p *Player, err error) {
	p = &Player{}
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Close()

	p.mixer = beep.Mixer{}

	numBytes := bufferSize * 4
	p.samples = make([][2]float64, bufferSize)
	p.buf = make([]byte, numBytes)

	p.context, err = oto.NewContext(int(sampleRate), 2, 2, numBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize speaker")
	}
	log.Print("before NewPlayer")
	p.player = p.context.NewPlayer()
	log.Print("before NewPlayer")

	p.done = make(chan struct{})
	log.Print("make channel")

	go func() {
		for {
			select {
			default:
				p.update()
			case <-p.done:
				return
			}
		}
	}()

	return p, nil
}

// Close closes the playback and the driver. In most cases, there is certainly no need to call Close
// even when the program doesn't play anymore, because in properly set systems, the default mixer
// handles multiple concurrent processes. It's only when the default device is not a virtual but hardware
// device, that you'll probably want to manually manage the device from your application.
func (p *Player) Close() {
	if p.player != nil {
		if p.done != nil {
			p.done <- struct{}{}
			p.done = nil
		}
		p.player.Close()
		p.context.Close()
		p.player = nil
	}
}

// Lock locks the speaker. While locked, speaker won't pull new data from the playing Stramers. Lock
// if you want to modify any currently playing Streamers to avoid race conditions.
//
// Always lock speaker for as little time as possible, to avoid playback glitches.
func (p *Player) Lock() {
	p.mu.Lock()
}

// Unlock unlocks the speaker. Call after modifying any currently playing Streamer.
func (p *Player) Unlock() {
	p.mu.Unlock()
}

// Play starts playing all provided Streamers through the speaker.
func (p *Player) Play(s ...beep.Streamer) {
	p.mu.Lock()
	p.mixer.Add(s...)
	p.mu.Unlock()
}

// Clear removes all currently playing Streamers from the speaker.
func (p *Player) Clear() {
	p.mu.Lock()
	p.mixer.Clear()
	p.mu.Unlock()
}

// update pulls new data from the playing Streamers and sends it to the speaker. Blocks until the
// data is sent and started playing.
func (p *Player) update() {
	p.mu.Lock()
	p.mixer.Stream(p.samples)
	p.mu.Unlock()

	// buf := p.buf
	for i := range p.samples {
		for c := range p.samples[i] {
			val := p.samples[i][c]
			if val < -1 {
				val = -1
			}
			if val > +1 {
				val = +1
			}
			valInt16 := int16(val * (1<<15 - 1))
			low := byte(valInt16)
			high := byte(valInt16 >> 8)
			p.buf[i*4+c*2+0] = low
			p.buf[i*4+c*2+1] = high
		}
	}

	p.player.Write(p.buf)
}
