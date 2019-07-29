// This package contains simple command line animations
package animator

import (
	"fmt"
	"time"

	"../color"
)

// Wait determines the base waiting time for any animator.
// Any speed will be relative to this value.
var Wait time.Duration

// A is just a sample const
const A = 10

// Animator is what animates the given input
type Animator struct {
	prefix string
	frames []string
	final  string
	speed  int64
	stop   bool
	ch     chan bool
	col    *color.Color
}

// New creates an animator with specified values.
// 'prefix' is printed before each frame.
// 'frames' are the different frames, e.g.: ".  ", ".. ", "...".
// The frames are played in order.
// 'final' is what should be printed when the frames have played, it will
// overwrite the previous frames if present.
// 'speed' is how fast it should be played, lower number = higher speed.
// If speed=1 the animation will update every 10ms.
func New(prefix string, frames []string, final string, speed int64) *Animator {
	a := &Animator{}
	a.prefix = prefix
	a.frames = frames
	a.final = final
	a.speed = speed
	a.ch = make(chan bool)
	return a
}

// Animate runs the animations indefinitely and will only stop if the Stop function is called.
func (a *Animator) Animate() {
	go a.animate()
}

// SetColor adds a color to the animator which will be used for all animation.
func (a *Animator) SetColor(c *color.Color) {
	a.col = c
}

func (a *Animator) animate() {
	l := len(a.frames)
	w, _ := time.ParseDuration(fmt.Sprintf("%dns", Wait.Nanoseconds()*a.speed))
	go a.wait()
	for !a.stop {
		for i := 0; i < l && !a.stop; i++ {
			a.print("\r" + a.prefix + a.frames[i])
			time.Sleep(w)
		}
	}
	if a.final != "" {
		a.print("\r" + a.final + "\n")
	} else {
		a.print("\n")
	}
	close(a.ch)
}

func (a *Animator) print(s string) {
	if a.col != nil {
		s = a.col.Apply(s)
	}
	fmt.Print(s)
}

func (a *Animator) wait() {
	defer func() {
		recover()
		a.stop = true
	}()
	_ = <-a.ch
	panic("Times up!")
}

// Stop stops the ongoing animations.
func (a *Animator) Stop() {
	a.ch <- true
	<-a.ch
}

func init() {
	Wait, _ = time.ParseDuration("10ms")
}
