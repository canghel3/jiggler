package main

import (
	"flag"
	"math/rand/v2"
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	defaultJiggleTime      = 10
	defaultJiggleFrequency = 10
)

func main() {
	f := flag.Uint("f", defaultJiggleFrequency, "jiggle frequency in seconds")
	t := flag.Uint("t", defaultJiggleTime, "jiggle time in minutes")

	flag.Parse()

	jiggler := newJiggler(*f, *t)
	if err := jiggler.run(); err != nil {
		panic(err)
	}
}

type Jiggler struct {
	screenHeight int
	screenWidth  int
	frequency    uint
	time         uint
}

func newJiggler(f, t uint) *Jiggler {
	w, h := robotgo.GetScreenSize()
	return &Jiggler{
		screenHeight: h,
		screenWidth:  w,
		frequency:    f,
		time:         t,
	}
}

func (j *Jiggler) run() error {
	ticker := time.NewTicker(time.Duration(j.frequency) * time.Second)
	timer := time.NewTimer(time.Duration(j.time) * time.Minute)

	for {
		select {
		case <-timer.C:
			return nil
		case <-ticker.C:
			if err := j.jiggle(); err != nil {
				return err
			}
		}
	}
}

func (j *Jiggler) jiggle() error {
	w := rand.N(j.screenWidth)
	h := rand.N(j.screenHeight)
	robotgo.MoveSmooth(w, h, 10.0, 10.0)
	return nil
}
