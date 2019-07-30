package spec

import (
	"fmt"
)

func log(a ...interface{}) {
  fmt.Fprintln(out, a...)
}

func vlog(i interface{}, prefix ...interface{}) {
  // Colored output doesn't work well with vim
	// fmt.Fprintf(out, "\x1b[38;2;%d;%d;%dm", 0xA0, 0xA0, 0x10)
	fmt.Fprint(out, prefix...)
  // fmt.Fprint(out, "\x1b[0m")
	fmt.Fprintf(out, "%#v\n", i)
}
