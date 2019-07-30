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
	fmt.Printf("\x1b[38;2;%d;%d;%dm", 0xA0, 0xA0, 0x10)
	fmt.Print(prefix...)
  fmt.Print("\x1b[0m")
	fmt.Printf("%#v\n", i)
}
