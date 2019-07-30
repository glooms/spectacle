package decls

type AType struct {
	a, b, c int
}

type AIntr interface {
	a() bool
	b(bool) bool
	c() func(bool) bool
}

type (
	Impl struct {
		val bool
	}

	Foo struct {
		a, b, c int
		d, x    bool
		A       AType
	}
)

type BType AType

type EmptyI interface{}
type EmptyS struct{}
