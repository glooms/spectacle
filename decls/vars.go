package decls

var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
	j       int
	u, v, s = 2.0, 3.0, "bar"
)
// var re, im, a = complexSqr(-1), 0
var entries = map[string]int{"Klas": 1}

var m1, m2, m3 map[bool]bool

// var _, found = entries["Klas"] // map lookup; only interested in "found"
// var b, c, d = entries["B"], entries["C"], entries["D"]


// var fa, fb, fc func (bool) int
