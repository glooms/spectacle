package spec

import (
  "fmt"
	"go/ast"
	"go/token"
  "strings"
)

type Package struct {
	Name   string

	Consts []*Value
	Vars   []*Value
	Types  []*ast.TypeSpec
	Funcs  []*ast.FuncDecl
}

type Value struct {
	Names     []string
  Types     []string
  Decl     *ast.ValueSpec
  Values    []string
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

func New(pkg *ast.Package) *Package {
	p := &Package{}

	p.Name = pkg.Name

	p.Consts = []*Value{}
	p.Vars = []*Value{}
	p.Types = []*ast.TypeSpec{}
	p.Funcs = []*ast.FuncDecl{}


	p.read(pkg)

	return p
}

func (p *Package) String() string {
	var str string
  defer func() {
    if r := recover(); r != nil {
      fmt.Println(str)
      panic(r)
    }
  }()
	// internal cannot be imported.
	if p.Name == "internal" {
		return ""
	}
	str += "package name: " + p.Name + "\n"
	str += "consts:\n"
	for _, c := range p.Consts {
    for i, name := range c.Names {
		  str += " " + name + " " + c.Types[i] + "\n"
    }
	}
	str += "vars:\n"
	for _, v := range p.Vars {
    for i, name := range v.Names {
		  str += " " + name + " " + v.Types[i] + "\n"
    }
	}
	return str
}

func (p *Package) read(pkg *ast.Package) {
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			p.readDecl(decl)
		}
	}
}

func (p *Package) readDecl(d ast.Decl) {
	// Ignoring imports right now, will be relevant for determining types.
	switch d.(type) {
	case *ast.GenDecl:
		d := d.(*ast.GenDecl)
		switch d.Tok {
		case token.VAR:

      for _, spec := range d.Specs {
        val := spec.(*ast.ValueSpec)
        p.Vars = append(p.Vars , readValue(val))
      }

		case token.CONST:

      for _, spec := range d.Specs {
        val := spec.(*ast.ValueSpec)
        p.Consts = append(p.Consts , readValue(val))
      }

		case token.TYPE:
			p.addTypes(d)
		}
	case *ast.FuncDecl:
		p.addFunc(d.(*ast.FuncDecl))
	}
}

func readValue(val *ast.ValueSpec) *Value {
  names := make([]string, len(val.Names))
  for i, id := range val.Names {
    names[i] = id.Name
  }
  fmt.Println(names)
  var types []string
  if val.Type != nil {
    typ := readExpr(val.Type)[0]
    types = make([]string, len(val.Names))
    for i := 0; i < len(types); i++ {
      types[i] = typ
    }
  } else {
    types = make([]string, 0, len(val.Names))
    for _, v := range val.Values {
      types = append(types, readExpr(v)...)
    }
  }
  v := &Value{
    Names: names,
    Types: types,
    Decl: val,
  }
  return v
}

func readExpr(e ast.Expr) []string {
	// TODO, lookup for more difficult value expressions
  fmt.Printf("lookup: %#v\n", e)
  types := []string{}

	switch e.(type) {
	case *ast.Ident:
    id := e.(*ast.Ident)
    if id.Obj == nil {
      types = append(types, id.Name)
    } else {
      readObj(id.Obj)
    }

  case *ast.BasicLit:
    lit := e.(*ast.BasicLit)
		types = append(types, strings.ToLower(lit.Kind.String()))

  case *ast.MapType:
    typ := e.(*ast.MapType)
    s := "map[" + readExpr(typ.Key)[0] + "]" + readExpr(typ.Value)[0]
    types = append(types, s)

  case *ast.FuncType:
    //typ := e.(*ast.FuncType)

	case *ast.CompositeLit:
		lit := e.(*ast.CompositeLit)
    types = readExpr(lit.Type)

	case *ast.StarExpr:
		types = readExpr(e.(*ast.StarExpr).X)

	case *ast.SelectorExpr:
		sel := e.(*ast.SelectorExpr)
    s := readExpr(sel.X)[0] + "." + readExpr(sel.Sel)[0]
		types = append(types, s)

	case *ast.CallExpr: // Function call, look at functions out params
    // Need to read decl
    //types = s.readExpr(cexpr.Fun)

	case *ast.IndexExpr: // Index expression of map, slice, array etc
    // Need to read decl
    types = append(types, "<none>", "<none>")
	}
  return types
}

func readObj(o *ast.Object) []string {
  // Decl related. Not yet implemented
  fmt.Printf("NYI: %#v\n", o)
  switch o.Kind {
  case ast.Bad:
  case ast.Pkg:
  case ast.Con:
  case ast.Typ:
  case ast.Var:
    if _, ok := o.Decl.(*ast.ValueSpec); ok {
      return nil
    }
  case ast.Fun:
    if _, ok := o.Decl.(*ast.FuncDecl); ok {
      return nil
    }
  case ast.Lbl:
  }
  return nil
}


func (p *Package) addTypes(d *ast.GenDecl) {
	for _, spec := range d.Specs {
		typ := spec.(*ast.TypeSpec)
		//id := typ.Name
		p.Types = append(p.Types, typ)
	}
}

func (p *Package) addFunc(d *ast.FuncDecl) {
	//id := d.Name
  p.Funcs = append(p.Funcs, d)
}
