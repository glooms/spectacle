package spec

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"sort"
  "strings"
)

type Spec struct {
	Name   string

	Consts map[string]*ast.ValueSpec
	Vars   map[string]*ast.ValueSpec
	Types  map[string]*ast.TypeSpec
	Funcs  map[string]*ast.FuncDecl
}

type Var struct {
	Name     string
	Type     string
	Comment  string
	complete bool
}

type Const struct {
	Name    string
	Type    string
	Comment string
	Value   string
}

type Type struct {
	Name    string
	Type    string
	Comment string
	Fields  []string // For struct
	Methods []*Func  // For interface, and maybe for structs with receiver funcs?
}

type Func struct {
	Name        string
	ParamNames  []string
	ParamTypes  []string
	ReturnNames []string // Can be nil
	ReturnTypes []string
	Receiver    string
}

func New(pkg *ast.Package) *Spec {
	s := &Spec{}

	s.Name = pkg.Name

	s.Consts = map[string]*ast.ValueSpec{}
	s.Vars = map[string]*ast.ValueSpec{}
	s.Types = map[string]*ast.TypeSpec{}
	s.Funcs = map[string]*ast.FuncDecl{}

	s.read(pkg)

	return s
}

func (s *Spec) String() string {
	var str string
	// internal cannot be imports are not allow externally.
	if s.Name == "internal" {
		return ""
	}
	str += "package name: " + s.Name + "\n"
	types := sillySort(reflect.ValueOf(s.Types))
	if len(types) > 0 {
		str += "types:\n"
	}
	for _, k := range types {
		str += " " + k + "\n"
	}
	consts := sillySort(reflect.ValueOf(s.Consts))
	if len(consts) > 0 {
		str += "consts:\n"
	}
	for _, k := range consts {
		str += " " + k + "\n"
	}
  str += s.docVars()
	funcs := sillySort(reflect.ValueOf(s.Funcs))
	if len(funcs) > 0 {
		str += "funcs:\n"
	}
	for _, k := range funcs {
		str += " " + k + "\n"
	}
	return str
}

func (s *Spec) docVars() string {
  str := ""
	if len(s.Vars) > 0 {
		str += "vars:\n"
	}
  types := map[string]string{}
  visited := map[*ast.ValueSpec]bool{}
	for _, v := range s.Vars {
    if visited[v] {
      continue
    }
    visited[v] = true
    t := s.reprVar(v)
    for i, id := range v.Names {
      fmt.Println(id.Name)
      types[id.Name] = t[i]
    }
	}
  vars := sillySort(reflect.ValueOf(s.Vars))
  for _, k := range vars {
    str += "  " + k + " " + types[k] + "\n"
  }
  return str
}

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

func (s *Spec) read(pkg *ast.Package) {
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			s.add(decl)
		}
	}
}

func (s *Spec) add(d ast.Decl) {
	// Ignoring imports right now, will be relevant for determining types.
	switch d.(type) {
	case *ast.GenDecl:
		d := d.(*ast.GenDecl)
		switch d.Tok {
		case token.VAR:
			s.addVar(d)
		case token.CONST:
			s.addConst(d)
		case token.TYPE:
			s.addTypes(d)
		}
	case *ast.FuncDecl:
		s.addFunc(d.(*ast.FuncDecl))
	}
}

func commentText(cg *ast.CommentGroup) string {
	var s string
	for _, c := range cg.List {
		s += c.Text + " "
	}
	return s
}

func (s *Spec) addVar(d *ast.GenDecl) {
	for _, spec := range d.Specs {
		val := spec.(*ast.ValueSpec)
		for _, id := range val.Names {
			s.Vars[id.Name] = val
		}
	}
}

func (s *Spec) reprVar(val *ast.ValueSpec) []string {
  hasType := val.Type != nil
  var types []string
  if hasType {
    typ := resolve(val.Type)
    types = make([]string, len(val.Names))
    for i := 0; i < len(types); i++ {
      types[i] = typ
    }
  } else {
    types = make([]string, 0, len(val.Names))
    for _, v := range val.Values {
      types = append(types, typeLookup(v)...)
    }
  }
  fmt.Printf("%+v\n", types)
  return types
}

func typeLookup(e ast.Expr) []string {
	// TODO, lookup for more difficult value expressions
  fmt.Printf("lookup: %#v\n", e)
  types := []string{}
	switch e.(type) {
	case *ast.Ident:
    types = append(types, resolve(e))
  case *ast.MapType:
    typ := e.(*ast.MapType)
    s := "map[" + resolve(typ.Key) + "]" + resolve(typ.Value)
    types = append(types, s)
	case *ast.BasicLit:
		types = append(types, resolve(e))
	case *ast.CompositeLit:
		lit := e.(*ast.CompositeLit)
    types = typeLookup(lit.Type)
	case *ast.StarExpr:
		types = typeLookup(e.(*ast.StarExpr).X)
	case *ast.SelectorExpr:
		sel := e.(*ast.SelectorExpr)
    s := resolve(sel.X) + "." + resolve(sel.Sel)
		types = append(types, s)
	case *ast.CallExpr: // Function call, look at functions out params
    cexpr := e.(*ast.CallExpr)
    types = typeLookup(cexpr.Fun)
	case *ast.IndexExpr: // Index expression of map, slice, array etc
    idx := e.(*ast.IndexExpr)
    types = typeLookup(idx.X)
	}
  fmt.Printf("%+v\n", types)
  return types
}

func resolve(e ast.Expr) string {
	fmt.Printf("resolve %#v\n", e)
	switch e.(type) {
	case *ast.Ident:
    id := e.(*ast.Ident)
    if id.Obj == nil {
		  return id.Name
    } else {
      fmt.Printf("%#v\n", id.Obj)
    }
	case *ast.BasicLit:
    lit := e.(*ast.BasicLit)
		return strings.ToLower(lit.Kind.String())
  }
  return ""
}

func (s *Spec) addConst(d *ast.GenDecl) {
	for _, spec := range d.Specs {
		val := spec.(*ast.ValueSpec)
		for _, id := range val.Names {
			s.Consts[id.Name] = val
		}
	}
}

func (s *Spec) addTypes(d *ast.GenDecl) {
	for _, spec := range d.Specs {
		typ := spec.(*ast.TypeSpec)
		id := typ.Name
		s.Types[id.Name] = typ
	}
}

func (s *Spec) addFunc(d *ast.FuncDecl) {
	id := d.Name
	s.Funcs[id.Name] = d
}
