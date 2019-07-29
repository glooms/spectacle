package decls

const a = 0
const (
  b, c = 1, 2
  d, e = 1.2, "hello"
)

const f = false
const g bool = true
const h = g
const (
  A int32 = iota
  B
  C
)
