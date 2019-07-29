package color

import (
	"encoding/hex"
	"fmt"
	"os"
)

type Color struct {
	R byte
	G byte
	B byte
}

func New(col string) *Color {
	c := &Color{}
	err := c.init([]byte(col))
	check(err)
	return c
}

func (c *Color) init(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("%s is not a valid color!", string(b))
	}
	if len(b) != 7 {
		return fmt.Errorf("%s is not a valid color!", string(b))
	}
	if rune(b[0]) != '#' {
		return fmt.Errorf("colors start with a '#' symbol")
	}
	c.R = dec(b[1:3])
	c.G = dec(b[3:5])
	c.B = dec(b[5:7])
	return nil
}

func (c *Color) String() string {
	s := fmt.Sprintf("#%X%X%X", c.R, c.G, c.B)
	s += c.color(" \u25A0\u25A0\u25A0\u25A0\u25A0\u25A0\u25A0")
	return s
}

func (c *Color) Apply(str string) string {
	return c.color(str)
}

func (c *Color) color(str string) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", c.R, c.G, c.B, str)
}

func dec(src []byte) byte {
	dst := make([]byte, 1)
	_, err := hex.Decode(dst, src)
	check(err)
	return dst[0]
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
