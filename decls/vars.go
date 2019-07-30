package decls

import "fmt"

var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
	j       int
	u, v, s = 2.0, 3.0, "bar"
)

var re, im = complexSqr(-1)
var entries = map[string]int{"Klas": 1}

var m1, m2, m3 map[bool]bool

var _, found = entries["Klas"] // map lookup; only interested in "found"
var _b, _c, _d = entries["B"], entries["C"], entries["D"]

var fa, fb, fc func(bool) int

var edge = func(b bool) bool {
	return b
}(true)

var str = fmt.Sprint("Hello!")

var foo = Foo{}
var sel = foo.A.a
