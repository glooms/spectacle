package spec

import (
  "errors"
  "go/ast"
  "go/importer"
	"go/parser"
  "go/token"
	"go/types"
  "sort"
)

// Spec represents a specification of a go package.
// It is usually built with the Build function but can also
// be built directly with New.
type Spec struct {
	Name string

	Consts []string
	Vars   []string
	Types  []string
	Funcs  []string

  Objects map[string]types.Object

  sorted bool
}

// New builds a Spec based on the contents of a directory.
//
// The function will be successful under two conditions:
//
//   1. There are .go files in the directory.
//   2. All .go files belong to the same package.
//
func New(path string) (*Spec, error) {
  var spec Spec

  spec.Consts = []string{}
  spec.Vars   = []string{}
  spec.Types  = []string{}
  spec.Funcs  = []string{}

  spec.Objects = map[string]types.Object{}

  err := spec.build(path)
  if err != nil {
    return nil, err
  }

  return &spec, nil
}

// build is the workhorse that builds a Spec.
func (s *Spec) build(path string) error {
  fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)

  if err != nil {
    return err
  }

  if len(pkgs) > 1 {
    return errors.New("found multiple packages")
  }

  if len(pkgs) == 0 {
    return errors.New("found no package")
  }

  var pkg *ast.Package
  for _, v := range pkgs {
    pkg = v
  }

  s.Name = pkg.Name

  conf := &types.Config{Importer: importer.Default()}

  files := make([]*ast.File, len(pkg.Files))
  i := 0
  for _, file := range pkg.Files {
    files[i] = file
    i++
  }

  p, err := conf.Check(path, fset, files, nil)
  if err != nil {
    return err
  }

  scope := p.Scope()
  for _, name := range scope.Names() {
    o := scope.Lookup(name)
    s.add(o)
  }

  return nil
}


// String is a template implementation of how Spec could be represented.
// Note that it contains newlines which might be frowned upon.
func (s *Spec) String() string {
  var str string

  if !s.sorted {
    sort.Strings(s.Consts)
    sort.Strings(s.Vars)
    sort.Strings(s.Funcs)
    sort.Strings(s.Types)
    s.sorted = true
  }

  str += "package: " + s.Name + "\n"

  str += "consts:\n"
  for _, name := range s.Consts {
    c := s.Objects[name].(*types.Const)
    typ := c.Type().String()
    val := c.Val().String()
    str += "  " + c.Name() + " " + typ + " " + val + "\n"
  }

  str += "types:\n"
  for _, name := range s.Types {
    t := s.Objects[name].(*types.TypeName)
    typ := t.Type().Underlying().String()
    str += "  " + t.Name() + " " + typ + "\n"
  }

  str += "vars:\n"
  for _, name := range s.Vars {
    v := s.Objects[name].(*types.Var)
    typ := v.Type().String()
    str += "  " + v.Name() + " " + typ + "\n"
  }

  str += "funcs:\n"
  for _, name := range s.Funcs {
    f := s.Objects[name].(*types.Func)
    typ := f.Type().String()
    str += "  " + f.Name() + " " + typ + "\n"
  }

  return str
}


// add adds an Object to the Spec. Objects other than
// Funcs, Consts, Vars or Types are ignored.
func (s *Spec) add(o types.Object) {

  switch o.(type) {

  case *types.Func:
    s.Funcs = append(s.Funcs, o.Name())

  case *types.Const:
    s.Consts = append(s.Consts, o.Name())

  case *types.Var:
    s.Vars = append(s.Vars, o.Name())

  case *types.TypeName:
    s.Types = append(s.Types, o.Name())

  default:
    return
  }
  s.Objects[o.Name()] = o

  s.sorted = false
}
