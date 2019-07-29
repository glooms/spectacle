package decls

var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
	j       int
	u, v, s = 2.0, 3.0, "bar"
)
var re, im, a = complexSqr(-1), 0
var entries = map[string]int{"Klas": 1}

//var _, found = entries["Klas"] // map lookup; only interested in "found"

func complexSqr(c complex64) (float32, float32) {
	return real(c), imag(c)
}
