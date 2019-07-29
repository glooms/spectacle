// This package contains an implementation of linked data structures.
// For now only a linked map.
package linked

import "fmt"

type Map struct {
	count uint64
	order map[string]uint64
	m     map[string]interface{}
}

func NewMap() *Map {
	return &Map{
		order: map[string]uint64{},
		m:     map[string]interface{}{},
	}
}

func (m *Map) Put(k string, v interface{}) interface{} {
	if _, ok := m.order[k]; ok {
		old := m.m[k]
		m.m[k] = v
		return old
	}
	m.m[k] = v
	m.order[k] = m.count
	m.count++
	return nil
}

func (m *Map) Get(k string) (v interface{}) {
	v, ok := m.m[k]
	if !ok {
		return nil
	}
	return v
}

func (m *Map) Keys() []string {
	a := make([]string, m.count-1)
	for k, i := range m.order {
		a[i] = k
	}
	return a
}

func (m *Map) String() string {
	l := m.count
	if l == 0 {
		return ""
	}
	strs := make([]string, l)
	for k, i := range m.order {
		v := m.m[k]
		strs[i] = "\"" + k + "\"" + ": "
		switch v.(type) {
		case string:
			strs[i] += "\"" + v.(string) + "\""
		default:
			strs[i] += fmt.Sprint(m.m[k])
		}
	}
	repr := "{"
	for _, s := range strs[:l-1] {
		repr += s + ", "
	}
	repr += strs[l-1] + "}"
	return repr
}
