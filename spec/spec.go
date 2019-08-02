package spec

import (
  "go/ast"
  "go/importer"
  "go/token"
	"go/types"
  "sort"
)

type Spec struct {
	Name string

	Consts []string
	Vars   []string
	Types  []string
	Funcs  []string

  Objects map[string]types.Object

  sorted bool
}

func New(fset *token.FileSet, pkg *ast.Package) *Spec {
  var spec Spec
  spec.Name = pkg.Name

  spec.Consts = []string{}
  spec.Vars   = []string{}
  spec.Types  = []string{}
  spec.Funcs  = []string{}

  spec.Objects = map[string]types.Object{}

  spec.build(fset, pkg)

  return &spec
}

func (s *Spec) build(fset *token.FileSet, pkg *ast.Package) error {
  path := pkg.Name
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
    s.addObj(o)
  }
  return nil
}


func (s *Spec) String() string {
  if !s.sorted {
    sort.Strings(s.Consts)
    sort.Strings(s.Vars)
    sort.Strings(s.Funcs)
    sort.Strings(s.Types)
    s.sorted = true
  }
  var str string

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

func (s *Spec) addObj(o types.Object) {

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

func TypeString(t types.Type) string {
  var s string
  vdebug(t, "type: ")

start:
  switch t.(type) {
  case *types.Basic:
  case *types.Pointer:
  case *types.Array:
  case *types.Slice:
  case *types.Map:
  case *types.Chan:
  case *types.Struct:
  case *types.Tuple:
  case *types.Signature:
  case *types.Named:
    t = t.Underlying()
    goto start
  case *types.Interface:
    t := t.(*types.Interface)
    for i := 0; i < t.NumMethods(); i++ {
      s += t.Method(i).String() + " "
    }
  default:
    t = t.Underlying()
    goto start
  }

  return s
}
