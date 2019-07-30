package spec

import (
	"fmt"
	"reflect"
	"sort"
)

// sillySort is admittedly a silly way to sort the keys of map.
// But, it exists since the maps we want to sort cannot be converted
// to map[string]interface{} and I refuse to write one sort per map-type.
func sillySort(m reflect.Value) []string {
	if m.Kind() != reflect.Map {
		return nil
	}
	sorted := make([]string, m.Len())
	filtered := 0
	for i, k := range m.MapKeys() {
		s := k.String()
		// if s != "" && unicode.IsUpper(rune(s[0])) {
		if true {
			sorted[i] = s
		} else {
			filtered++
		}
	}
	sort.Strings(sorted)
	return sorted[filtered:]
}

func vprint(i interface{}, prefix ...interface{}) {
	fmt.Print(prefix...)
	fmt.Printf("%#v\n", i)
}
