package spec

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/printer"
	"go/token"
  "io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var out *os.File
var debug_out *os.File

var specs map[string]*Spec

// Build goes through all subdirectories of the specifed
// directory and create a new Spec for each one.
//
// A map containing all found Specs is returned
//
// If exported is true, only exported consts, types, etc. is
// included in each specification.
//
func Build(root string) map[string]*Spec {
  specs = map[string]*Spec{}
	err := filepath.Walk(root, walker)
	check(err)
  return specs
}

func init() {
  if _, err := os.Stat("./log"); os.IsNotExist(err) {
    os.Mkdir("./log", 0755)
  }
}

// filter is used to filter out go test files.
func filter(fi os.FileInfo) bool {
	return !strings.Contains(fi.Name(), "test")
}

// walker is the walking function used for going through all subdirectories
// of a specified directory.
func walker(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return nil
	}
  // Don't go through hidden files, probably not go code in there.
  if path[0] == '.' {
    return nil
  }
  files, err := ioutil.ReadDir(path)
  if err != nil {
    return err
  }
  // Check if there are any .go files in the directory.
  var any bool
  for _, f := range files {
    if strings.Contains(f.Name(), ".go") {
      any = true
      break
    }
  }
  if !any {
    return nil
  }

  createLogs(path)
  defer closeLogs()

  spec, _ := New(path)
  if spec == nil {
    return nil
  }
  fmt.Println("Found: '" + spec.Name + "'")
  log(spec.String())
  specs[spec.Name] = spec

	return nil
}

func createLogs(name string) {
  var err error
  out, err = os.Create("./log/" + name + ".log")
  check(err)
	debug_out, err = os.Create("./log/" + name + ".debug")
  check(err)
}

func closeLogs() {
  out.Close()
  debug_out.Close()
}

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func log(a ...interface{}) {
	fmt.Fprintln(out, a...)
}

func vlog(i interface{}, prefix ...interface{}) {
	fmt.Fprint(out, prefix...)
	fmt.Fprintf(out, "%#v\n", i)
}

func debug(a ...interface{}) {
	fmt.Fprintln(debug_out, a...)
}

func vdebug(i interface{}, prefix ...interface{}) {
	fmt.Fprint(debug_out, prefix...)
	fmt.Fprintf(debug_out, "%#v\n", i)
}

// docPrint pretty prints all decls using "go/doc" and "go/printer".
// Use it as a reference when needed.
func docPrint(pkg *ast.Package, fset *token.FileSet, path string) {
	p := doc.New(pkg, path, doc.AllDecls)
	fmt.Println("doc:", p.Doc)
	fmt.Println("name:", p.Name)
	fmt.Println("imports:")
	for _, imp := range p.Imports {
		fmt.Println("  " + imp)
	}
	fmt.Println("files:")
	for _, file := range p.Filenames {
		fmt.Println("  " + file)
	}
	fmt.Println("consts:")
	for _, c := range p.Consts {
		printer.Fprint(os.Stdout, fset, c.Decl)
		fmt.Println()
	}
	fmt.Println("types:")
	for _, t := range p.Types {
		printer.Fprint(os.Stdout, fset, t.Decl)
		fmt.Println()
	}
	fmt.Println("vars:")
	for _, v := range p.Vars {
		printer.Fprint(os.Stdout, fset, v.Decl)
		fmt.Println()
	}
	fmt.Println("funcs:")
	for _, f := range p.Funcs {
		printer.Fprint(os.Stdout, fset, f.Decl)
		fmt.Println()
	}
}
