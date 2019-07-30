package decls

func complexSqr(c complex64) (float32, float32) {
	return real(c), imag(c)
}

func (a *AType) do() (x, y, z int) {
	return
}

func (AType) ok() {
	return
}

func (i Impl) a() bool {
	return i.val
}

func (i Impl) b(a bool) bool {
	return i.a()
}

func (i Impl) c() func(bool) bool {
	return i.b
}
