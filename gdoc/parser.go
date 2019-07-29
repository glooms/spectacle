// gdoc = godoc doc. First draft at a API spec generator using go doc
// Doesn't handle multiple return types or imports
// Can generate a JSON object that resembles what's expected.
package gdoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"../linked"
)

var types = map[string]bool{}
var defs *linked.Map

func init() {
	defs = linked.NewMap()
}

func Generate(path string) {
	doc := run("go", "doc", "-all", path)
	fmt.Println(doc)
	parse(doc)
	buf := bytes.Buffer{}
	err := json.Indent(&buf, []byte(defs.String()), "", "  ")
	isFatal(err)
	fmt.Println(buf.String())
}

func addDef(def, name, typ string) {
	d := linked.NewMap()
	d.Put("def", def)
	if typ != "=" && typ != "" {
		d.Put("type", typ)
	}
	defs.Put(name, d)
}

func addFunc(strct, name, param, ret string) {
	d := linked.NewMap()
	d.Put("params", param)
	d.Put("returns", ret)
	if strct != "" {
		if s := strings.Split(strct, " ")[1]; s[0] == '*' {
			strct = s[1:]
		}
		m := defs.Get(strct)
		if m, ok := m.(*linked.Map); ok {
			m.Put(name, d)
		}
	} else {
		defs.Put(name, d)
	}
}

func parse(doc string) {
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		if line != "" {
			parseLine(line)
		}
	}
}

func parseLine(line string) {
	i := 0
	buf := []byte{}
	for ; i < len(line) && line[i] != ' '; i++ {
		buf = append(buf, line[i])
	}
	i++
	switch literal := string(buf); literal {
	case "const":
		fallthrough
	case "var":
		fallthrough
	case "type":
		buf = []byte{}
		for ; i < len(line) && line[i] != ' '; i++ {
			buf = append(buf, line[i])
		}
		i++
		name := string(buf)
		buf = []byte{}
		for ; i < len(line) && line[i] != ' '; i++ {
			buf = append(buf, line[i])
		}
		i++
		typ := string(buf)
		addDef(literal, name, typ)
	case "func":
		buf = []byte{}
		var strct, name, param, ret string
		for ; i < len(line); i++ {
			switch line[i] {
			case '(':
				i++
				if len(buf) != 0 {
					name = string(buf)
					buf = []byte{}
				}
				for ; line[i] != ')'; i++ {
					buf = append(buf, line[i])
				}
				i++
				if name == "" {
					strct = string(buf)
				} else {
					param = string(buf)
				}
				buf = []byte{}
			case ' ':
				continue
			default:
				buf = append(buf, line[i])
			}
		}
		if name == "" {
			name = string(buf)
		} else {
			ret = string(buf)
		}
		addFunc(strct, name, param, ret)
	}
}

func run(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()
	isFatal(err)
	return string(out)
}

func isFatal(e error) {
	if e != nil {
		code := 1
		switch e.(type) {
		case *exec.ExitError:
			err := e.(*exec.ExitError)
			fmt.Println(string(err.Stderr))
			code = err.ExitCode()
		default:
			fmt.Println(e)
		}
		os.Exit(code)
	}
}
