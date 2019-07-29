package main

import (
	"time"

	"./animator"
	"./color"
)

func example() {
	a := animator.New("Hello", []string{
		".        ",
		"..       ",
		"...      ",
		"....     ",
		".....    ",
		"......   ",
		".......  ",
		"........ ",
		".........",
	}, "Hello Universe!", 10)
	c := color.New("#FF0000")
	a.SetColor(c)
	a.Animate()
	dur, _ := time.ParseDuration("5s")
	time.Sleep(dur)
	a.Stop()
}
