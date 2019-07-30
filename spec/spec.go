package spec

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type Package struct {
	Name string

	Consts []*Value
	Vars   []*Value
	Types  []*Type
	Funcs  []*Func

  unread []*ast.Decl
}

type Value struct {
	Names  []string
	Decl   *ast.ValueSpec
	Types  []string
	Values []string
}

type Type struct {
	Name       string
	Decl       *ast.TypeSpec
	Type       string
	Comment    string
	FieldNames []string // For struct
	FieldTypes []string // For struct
	MethodNames     []string // For interface
	MethodTypes    []string // For interface
}

type Func struct {
	Name        string
	Decl        *ast.FuncDecl
	ParamNames  []string
	ParamTypes  []string
	ReturnNames []string // Can be nil
	ReturnTypes []string // Can be nil
	Receiver    string
}

var cur_pkg *Package

func New(pkg *ast.Package) *Package {
	p := &Package{}

	p.Name = pkg.Name

	p.Consts = []*Value{}
	p.Vars = []*Value{}
	p.Types = []*Type{}
	p.Funcs = []*Func{}

  p.unread = []*ast.Decl{}

	p.read(pkg)

  cur_pkg = p

	return p
}

func (p *Package) lookupFunc(name string) *Func {
  for _, f := range p.Funcs {
    if f.Name == name {
      return f
    }
  }
  return nil
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
			if i < len(c.Types) {
				str += "  " + name + " " + c.Types[i] + "\n"
			} else {
				str += "  " + name + " INCOMPLETE\n"
			}
		}
	}

  /* str += "types:\n"
  for _, t := range p.Types {
    str += "  type " + t.Name + " "
    switch t.Type {
    case "struct":

      str += t.Type + " {"
      i := 0
      for ; i < len(t.FieldNames); i++ {
        str += "\n    " + t.FieldNames[i] + " " + t.FieldTypes[i]
      }
      if i > 0 {
        str += "\n  "
      }
      str += "}"

    case "interface":

      str += t.Type + " {"
      i := 0
      for ; i < len(t.MethodNames); i++ {
        str += "\n    " + t.MethodNames[i] + " " + t.MethodTypes[i]
      }
      if i > 0 {
        str += "\n  "
      }
      str += "}"

    default:
      str += t.Type
    }
    str += "\n"
  } */

	str += "vars:\n"
	for _, v := range p.Vars {
		for i, name := range v.Names {
			if i < len(v.Types) {
				str += "  " + name + " " + v.Types[i] + "\n"
			} else {
				str += "  " + name + " INCOMPLETE\n"
			}
		}
	}

	/* str += "funcs:\n"
	for _, f := range p.Funcs {
		str += "  func "
		if f.Receiver != "" { // Maybe these should just be added to the receiver
			str += "(" + f.Receiver + ") "
		}
		str += f.Name + "("
		for i, par := range f.ParamNames {
			str += par + " " + f.ParamTypes[i]
			if i < len(f.ParamNames)-1 {
				str += ", "
			}
		}
		str += ")"
		if f.ReturnTypes == nil {
			str += "\n"
			continue
		}
		str += " "
		paren := len(f.ReturnTypes) > 1
		if paren {
			str += "("
		}
		for i, ret := range f.ReturnTypes {
			str += ret
			if i < len(f.ReturnTypes)-1 {
				str += ", "
			}
		}
		if paren {
			str += ")"
		}
		str += "\n"
	}
  */

	return str
}

func (p *Package) read(pkg *ast.Package) {
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			p.readDecl(decl)
		}
	}
}

// readDecl reads a Decl, which is an interface type, and determines what
// sort of declaration it is and calls the appropriate read-function.
func (p *Package) readDecl(d ast.Decl) {
	// Ignoring imports right now, will be relevant for determining types.
	switch d.(type) {

	case *ast.GenDecl:
		d := d.(*ast.GenDecl)
		switch d.Tok {
		case token.VAR:

			for _, spec := range d.Specs {
				val := spec.(*ast.ValueSpec)
				p.Vars = append(p.Vars, readValue(val))
			}

		case token.CONST:

			for _, spec := range d.Specs {
				val := spec.(*ast.ValueSpec)
				p.Consts = append(p.Consts, readValue(val))
			}

		case token.TYPE:

			for _, spec := range d.Specs {
				typ := spec.(*ast.TypeSpec)
				p.Types = append(p.Types, readType(typ))
			}
		}

	case *ast.FuncDecl:
		fun := d.(*ast.FuncDecl)
		p.Funcs = append(p.Funcs, readFunc(fun))
	}
}

// readValue reads a ValueSpec to map variables/constants
// to types and their initial value. (Mostly important in the
// constant case.)
//
// To be implemented:
//   1. assign initial values (if any)
//   2. handle references
//
func readValue(val *ast.ValueSpec) *Value {
	names := make([]string, len(val.Names))
	for i, id := range val.Names {
		names[i] = id.Name
	}

	var types []string
	if val.Type != nil {
		typ := readTypeExpr(val.Type)[0]
		types = make([]string, len(val.Names))
		for i := 0; i < len(types); i++ {
			types[i] = typ
		}
	} else {
		types = make([]string, 0, len(val.Names))
		for _, v := range val.Values {
			types = append(types, readTypeExpr(v)...)
		}
	}

	v := &Value{
		Names: names,
		Types: types,
		Decl:  val,
	}
	return v
}

// readFunc reads a FuncDecl and creates a Func.
// It extracts information about the function parameters and return
// types as well as whether or not it has a receiver (which would technically make it method).
//
// We might want to add methods to their receivers, but that can be done later too.
func readFunc(d *ast.FuncDecl) *Func {
	name := d.Name.Name

	parNames, parTypes := readFields(d.Type.Params)
	retNames, retTypes := readFields(d.Type.Results)
	var recv string

	if d.Recv != nil {
		_, recvTypes := readFields(d.Recv)
		recv = recvTypes[0]
	}

	f := &Func{
		Name:        name,
		Decl:        d,
		ParamNames:  parNames,
		ParamTypes:  parTypes,
		ReturnNames: retNames,
		ReturnTypes: retTypes,
		Receiver:    recv,
	}
	return f
}

// readType reads a TypeSpec to create a Type.
// A TypeSpec can declare a new struct or interface or reference some
// pre-existing basic type or previously imported/declared type.
//
// Struct declarations (StructType) will populated FieldNames/Types.
//
// Interface declarations (InterfaceType) will populated MethodNames/Types.
func readType(typ *ast.TypeSpec) *Type {

	var typeType string
	var fieldNames, fieldTypes []string
  var methodNames, methodTypes []string

	switch typ.Type.(type) { // What type is typ's Type's type?

	case *ast.StructType:
		typeType = "struct"
		strct := typ.Type.(*ast.StructType)
		fieldNames, fieldTypes = readFields(strct.Fields)

	case *ast.InterfaceType:
		typeType = "interface"
		intr := typ.Type.(*ast.InterfaceType)
		methodNames, methodTypes = readFields(intr.Methods)

	default:
		typeType = readTypeExpr(typ.Type)[0]
	}

	t := &Type{
		Name:       typ.Name.Name,
		Decl:       typ,
		Type:       typeType,
		FieldNames: fieldNames,
		FieldTypes: fieldTypes,
    MethodNames: methodNames,
    MethodTypes: methodTypes,
	}
	return t
}

// readFields reads a FieldList and returns any names and their corresponding types.
// If there are no names the type will still be returned as it must be set.
// For example a function like:
//
//   func a(b, c int) int
//
// would have two FieldLists, one for params and one for return values.
// The params would be two variables of type int and the return an unnamed variable
// of type int.
func readFields(fl *ast.FieldList) (names []string, types []string) {
	if fl == nil {
		return
	}
	names = []string{}
	types = []string{}
	for _, f := range fl.List {
		typ := readTypeExpr(f.Type)[0]
		if f.Names != nil {
			for _, n := range f.Names {
				names = append(names, n.Name)
				types = append(types, typ)
			}
		} else {
			types = append(types, typ)
		}
	}
	return
}

// readTypeExpr reads an expression and returns the associated types.
//
// A lot of cases remain to be implemented.
// Particularly any Expr that reference another declaration.
//
func readTypeExpr(e ast.Expr) []string {
	//vprint(e, "readTypeExpr: ")
	types := []string{}

	switch e.(type) {

	case *ast.Ident:
		id := e.(*ast.Ident)
		if id.Obj == nil {
			types = append(types, id.Name)
		} else {
			types = append(types, readObj(id.Obj)...)
		}

	case *ast.BasicLit:
		lit := e.(*ast.BasicLit)
		types = append(types, strings.ToLower(lit.Kind.String()))

	case *ast.MapType:
		typ := e.(*ast.MapType)
		s := "map[" + readTypeExpr(typ.Key)[0] + "]" + readTypeExpr(typ.Value)[0]
		types = append(types, s)

	case *ast.FuncType:
		// This should probably be refactored to a function
		typ := e.(*ast.FuncType)
		_, pars := readFields(typ.Params)
		_, rvs := readFields(typ.Results)
		s := "func(" + strings.Join(pars, ", ") + ") "
		paren := len(rvs) > 1
		if paren {
			s += "("
		}
		s += strings.Join(rvs, ", ")
		if paren {
			s += ")"
		}
		types = append(types, s)

	case *ast.CompositeLit:
		lit := e.(*ast.CompositeLit)
		types = readTypeExpr(lit.Type)

	case *ast.StarExpr:
		types = readTypeExpr(e.(*ast.StarExpr).X)

	case *ast.SelectorExpr:
		sel := e.(*ast.SelectorExpr)
		s := readTypeExpr(sel.X)[0] + "." + readTypeExpr(sel.Sel)[0]
		types = append(types, s)

  default:
    // Probably not a simple type expression, read it as a use expression.
    types = readUseExpr(e)
  }
	return types
}

func readUseExpr(e ast.Expr) []string {
  vprint(e, "readUseExpr: ")
  types := []string{}

  switch e.(type) {

  case *ast.IndexExpr:
    types = readUseExpr(e.(*ast.IndexExpr).X)

  case *ast.MapType:
    // This is the case where the parent was an IndexExpr. In this case,
    // dependent on how many variables are assigned the value we would
    // want to add a "bool" type to the list of types.
    // Example:
    //
		//   var m map[int]int
		//   var a     = m[0]
    //   var b, ok = m[0]
    //   var c, d  = m[0], m[0]
    //
    // For 'a' we simply want to return the found type.
    // For 'b', 'ok' we want to return the found type and a bool.
    // The last case shows why we can't always just add a bool to the
    // types as this would be an erroneous type for 'd'.
		typ := e.(*ast.MapType)
		types = readTypeExpr(typ.Value)

  case *ast.CallExpr:

    call := e.(*ast.CallExpr)

    switch call.Fun.(type) {
    case *ast.Ident:
      // Do lookup on name
      id := call.Fun.(*ast.Ident)
      types = cur_pkg.lookupFunc(id.Name).ReturnTypes

    case *ast.FuncLit:
      fun := call.Fun.(*ast.FuncLit)
      types = readUseExpr(fun.Type)
    }

	case *ast.FuncType:
		typ := e.(*ast.FuncType)
		_, rvs := readFields(typ.Results)
    types = append(types, rvs...)
  }
  return types
}

// readObj reads an Object which seem to be mostly used when variables
// reference other declarations and such. This remains to be implemented.
func readObj(o *ast.Object) []string {
	// Decl related. Not yet implemented
	vprint(o, "readObj: ")
	switch o.Kind {
	case ast.Bad:
	case ast.Pkg:
	case ast.Con:
	case ast.Typ:
		if t, ok := o.Decl.(*ast.TypeSpec); ok {
			return []string{t.Name.Name}
		}
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
