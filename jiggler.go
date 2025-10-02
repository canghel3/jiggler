package main

import (
	"flag"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	gohook "github.com/robotn/gohook"
)

const (
	defaultJiggleTime = 10
)

func main() {
	t := flag.Uint("t", defaultJiggleTime, "jiggle time in minutes")

	flag.Parse()

	jiggler := newJiggler(*t)
	if err := jiggler.run(); err != nil {
		panic(err)
	}
}

type Jiggler struct {
	screenHeight int
	screenWidth  int
	time         uint
}

func newJiggler(t uint) *Jiggler {
	w, h := robotgo.GetScreenSize()
	return &Jiggler{
		screenHeight: h,
		screenWidth:  w,
		time:         t,
	}
}

func (j *Jiggler) run() error {
	timer := time.NewTimer(time.Duration(j.time) * time.Minute)
	go j.stop()

	for {
		select {
		case <-timer.C:
			return nil
		default:
			if err := j.jiggleHumanLike(); err != nil {
				return err
			}
		}
	}
}

func (j *Jiggler) stop() {
	defer gohook.End()
	for event := range gohook.Start() {
		if event.Kind == gohook.KeyDown {
			switch event.Keycode {
			case gohook.Keycode[robotgo.KeyQ]:
				os.Exit(1)
			}
		}
	}
}

func (j *Jiggler) jiggle() error {
	w := rand.N(j.screenWidth)
	h := rand.N(j.screenHeight)
	s := time.Now()
	log.Printf("moving to %d %d\n", w, h)
	robotgo.MoveSmooth(w, h, 0.7, 1.1, 500)
	log.Printf("took %v to move\n", time.Since(s))
	return nil
}

func (j *Jiggler) jiggleHumanLike() error {
	startX, startY := robotgo.Location()
	destX := rand.N(j.screenWidth)
	destY := rand.N(j.screenHeight)

	log.Printf("starting from %d %d\nfinishing at %d %d\n", startX, startY, destX, destY)

	p0, p3 := point{startX, startY}, point{destX, destY}
	p1, p2 := getControlPointsInBbox(p0, p3)

	s := time.Now()

	steps := 100
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		curvePoint := cubicBezier(p0, p1, p2, p3, t)
		log.Printf("moving to %d %d\n", curvePoint.x, curvePoint.y)
		//determine move speed based on distance to cover. small distance in human-like behaviour means small speed, whereas large distance means higher speed
		robotgo.MoveSmooth(curvePoint.x, curvePoint.y, 0.01, 0.05, 1)
	}
	log.Printf("took %v to move\n", time.Since(s))
	return nil
}

type point struct {
	x, y int
}

func getControlPointsInBbox(p0, p3 point) (point, point) {
	minx := min(p0.x, p3.x)
	miny := min(p0.y, p3.y)
	maxx := max(p0.x, p3.x)
	maxy := min(p0.y, p3.y)

	boxWidth := maxx - minx
	boxHeight := maxy - miny

	return point{minx + rand.N(boxWidth+1), miny + rand.N(boxHeight+1)}, point{minx + rand.N(boxWidth+1), miny + rand.N(boxHeight+1)}
}

func distance(p0, p1 point) float64 {
	return math.Sqrt(float64((p0.x-p1.x)^2) + float64((p0.y-p1.y)^2))
}

func cubicBezier(p0, p1, p2, p3 point, t float64) point {
	t2 := t * t
	t3 := t2 * t
	mt1 := 1 - t
	mt2 := mt1 * mt1
	mt3 := mt2 * mt1

	//B(t) = (1-t)^3 * P0 + 3 * (1-t)^2 * t * P1 + 3 * (1-t) * t^2 * P2 + t^3 * P3; 0 <= t <= 1
	x := mt3*float64(p0.x) + 3*mt2*t*float64(p1.x) + 3*mt1*t2*float64(p2.x) + t3*float64(p3.x)
	y := mt3*float64(p0.y) + 3*mt2*t*float64(p1.y) + 3*mt1*t2*float64(p2.y) + t3*float64(p3.y)

	return point{
		x: int(x),
		y: int(y),
	}
}
