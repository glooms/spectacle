package spec

import (
	"go/ast"
	"go/token"
	"strings"
)

type Package struct {
	Name string

  Imports map[string]*ast.Object

	Consts []*Value
	Vars   []*Value
	Types  []*Type
	Funcs  []*Func
}

type Value struct {
	Names  []string
	Decl   *ast.ValueSpec
	Types  []string
	Values []string
}

type Type struct {
	Name        string
	Decl        *ast.TypeSpec
	Type        string
	Comment     string
	FieldNames  []string // For struct
	FieldTypes  []string // For struct
	MethodNames []string // For interface
	MethodTypes []string // For interface
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

func NewPkg(pkg *ast.Package) *Package {
	p := &Package{}

	p.Name = pkg.Name

	p.Consts = []*Value{}
	p.Vars = []*Value{}
	p.Types = []*Type{}
	p.Funcs = []*Func{}

  cur_pkg = p

  p.read(pkg)

	return p
}

var lookup = map[string]string {
  "false": "bool",
  "true": "bool",
}

func (p *Package) lookupFunc(name string) *Func {
	for _, f := range p.Funcs {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (p *Package) lookupType(name string) *Type {
  for _, t := range p.Types {
    if t.Name == name {
      return t
    }
  }
  return nil
}

func (p *Package) lookupValue(decl *ast.ValueSpec) *Value {
  for _, v := range p.Vars {
    if v.Decl == decl {
      return v
    }
  }
  for _, c := range p.Consts {
    if c.Decl == decl {
      return c
    }
  }
  return nil
}

func (p *Package) String() string {
	var str string
	defer func() {
		if r := recover(); r != nil {
			log(str)
			panic(r)
		}
	}()
	// internal cannot be imported.
	if p.Name == "internal" {
		return ""
	}
  warn := "\x1b[38;2;200;00;100m" + " INCOMPLETE" + "\x1b[0m"

	str += "package name: " + p.Name + "\n"

	str += "consts:\n"
	for _, c := range p.Consts {
		for i, name := range c.Names {
			str += "  " + name
			if i < len(c.Types) {
        str += " " + c.Types[i]
      }
      if i < len(c.Values) {
        str += " " + c.Values[i]
      }
      str += "\n"
		}
	}

	str += "types:\n"
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
	}

	str += "vars:\n"
	for _, v := range p.Vars {
		for i, name := range v.Names {
			if i < len(v.Types) {
				str += "  " + name + " " + resolveType(v.Types[i]) + "\n"
			} else {
				str += "  " + name + warn + "\n"
			}
		}
	}

	str += "funcs:\n"
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

	return str
}

func (p *Package) read(pkg *ast.Package) {
  p.Imports = pkg.Imports
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
				p.Vars = append(p.Vars, readValueSpec(val))
			}

		case token.CONST:

			for _, spec := range d.Specs {
				val := spec.(*ast.ValueSpec)
				p.Consts = append(p.Consts, readValueSpec(val))
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

// readValueSpec reads a ValueSpec to map variables/constants
// to types and their initial value. (Mostly important in the
// constant case.)
//
// To be implemented:
//   1. assign initial values (if any)
//   2. handle references
//
func readValueSpec(val *ast.ValueSpec) *Value {
	names := make([]string, len(val.Names))
	for i, id := range val.Names {
		names[i] = id.Name
	}

  log(names)

  // We should probably handle the two cases here differently.
  // If val.Type is set it means we have a value specification of the form:
  //
  //   1. var a, b, ... int
  //
  // which is quite easy to read as all variables have the same type.
  // The only things we have to look out for are any specified values in
  // case the specification is for a const.
  //
  // The other case however is when the type is implicit and val.Type is nil.
  // Could be for example one of:
  //
  //   2. var a, ok = m[0]
  //      var b, c  = foo.Foo, foo.Foo.foo()
  //
  // which are much trickier to determine.
  // We should probably use the fact that the type in the first case is limited
  // to just a few types of expressions.

	var types, values []string
  hasType := val.Type != nil
	if hasType {
		typ := readTypeExpr(val.Type)[0]
		types = make([]string, len(val.Names))
		for i := 0; i < len(types); i++ {
			types[i] = typ
		}
	}
  if val.Values != nil {
    multiReturn := false
    if !hasType {
	    types = make([]string, 0, len(val.Names))
      multiReturn = len(val.Values) < len(val.Names)
    }
    values = make([]string, 0, len(val.Names))
		for _, v := range val.Values {
      rtypes, rvalues := readValueExpr(v)
      if !hasType && rtypes != nil {
        if multiReturn {
          types = append(types, rtypes...)
        } else {
          types = append(types, rtypes[0])
        }
      }
      if rvalues != nil {
        values = append(values, rvalues...)
      }
		}
	}

	v := &Value{
		Names: names,
		Types: types,
		Decl:  val,
    Values: values,
	}
	return v
}

// resloveType takes a string representing a type such as:
//
//   foo.Foo.foo()
//
// and tries to resolve what type such a selection would return.
func resolveType(typstr string) string {
  parts := strings.Split(typstr, ".")
  if len(parts) > 1 {
    typ := cur_pkg.lookupType(parts[0])
    for _, part := range parts[1:] {
typeLoop:
      for i, name := range typ.FieldNames {
        if name == part {
          ftype := typ.FieldTypes[i]
          typ = cur_pkg.lookupType(ftype)
          if typ == nil {
            return ftype
          }
          break typeLoop
        }
      }
    }
  }
  return typstr
}

// readFunc reads a FuncDecl and creates a Func.
// It extracts information about the function parameters and return
// types as well as whether or not it has a receiver (which would technically make it method).
//
// We might want to add methods to their receivers, but that can be done later too.
func readFunc(d *ast.FuncDecl) *Func {
	name := d.Name.Name

  log(d.Name.Name)

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

  log(typ.Name.Name)

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
		Name:        typ.Name.Name,
		Decl:        typ,
		Type:        typeType,
		FieldNames:  fieldNames,
		FieldTypes:  fieldTypes,
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
// Value lookup remaings. Mostly relevant for constants but could also
// be nice for variables.
//
func readTypeExpr(e ast.Expr) []string {
	vlog(e, "type expr: ")

	switch e.(type) {

	case *ast.Ident:
		id := e.(*ast.Ident)
		if id.Obj == nil {
			return []string{id.Name}
		}
    return readTypeObj(id.Obj)

	case *ast.BasicLit:
		lit := e.(*ast.BasicLit)
		return []string{strings.ToLower(lit.Kind.String())}

	case *ast.MapType:
		typ := e.(*ast.MapType)
		s := "map[" + readTypeExpr(typ.Key)[0] + "]" + readTypeExpr(typ.Value)[0]
		return []string{s}

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
		return []string{s}

	case *ast.StarExpr:
		return readTypeExpr(e.(*ast.StarExpr).X)

	case *ast.SelectorExpr:
    // This case represents expressions of the form:
    //
    //   1. fmt.Print()    call of exposed function
    //   2. fmt.SomeConst  access of exposed variable/constant
    //   3. A.a            access of field in struct
    //   4. A.b()          call of method in struct/interface
    //
    // As you can see, they can be quite different.
		sel := e.(*ast.SelectorExpr)
    s := readTypeExpr(sel.X)[0] + "." + readTypeExpr(sel.Sel)[0]
		return []string{s}
	}
  return nil
}

func readValueExpr(e ast.Expr) (types, values []string) {
	vlog(e, "value expr: ")

	switch e.(type) {

  case *ast.Ident:
    id := e.(*ast.Ident)
    o := id.Obj
    if o == nil {
      types = []string{lookup[id.Name]}
      values = []string{id.Name}
      return types, values
    }
    return readValueObj(o)

  case *ast.BasicLit:
		lit := e.(*ast.BasicLit)
    types = []string{strings.ToLower(lit.Kind.String())}
    values = []string{lit.Value}
		return

	case *ast.IndexExpr:
    types, values = readValueExpr(e.(*ast.IndexExpr).X)
    // 'types' is list containing one elemet which is one of:
    //    map[T]T, []T, [x]T
    //
    types = []string{strings.Split(types[0], "]")[1]}
    types = append(types, "bool")
		return

  case *ast.CompositeLit:
    // Composite literal is a value expression. However, it is a
    // type expression followed by expressions within braces. We
    // can therefore lookup the type (but we will ignore the value for now).
		lit := e.(*ast.CompositeLit)
		return readTypeExpr(lit.Type), nil

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
    log("called map!")
		typ := e.(*ast.MapType)
		return readTypeExpr(typ.Value), nil

	case *ast.CallExpr:

		call := e.(*ast.CallExpr)

		switch call.Fun.(type) {

		case *ast.Ident:
			// Do lookup on name
			id := call.Fun.(*ast.Ident)
      f := cur_pkg.lookupFunc(id.Name)
      if f != nil {
        return f.ReturnTypes, nil
      }
      return nil, nil

    default:
      return readValueExpr(call.Fun)
		}

 case *ast.SelectorExpr:
   // Getting in to annoying territory.
   sel := e.(*ast.SelectorExpr)
   _, _ = readValueExpr(sel.X)
   _, _ = readValueExpr(sel.Sel)
   return nil, nil

  case *ast.FuncLit:
    lit := e.(*ast.FuncLit)
    return readValueExpr(lit.Type)

	case *ast.FuncType:
		typ := e.(*ast.FuncType)
    _, types := readFields(typ.Results)
		return types, nil
	}
	return nil, nil
}

// readTypeObj reads an Object which seem to be mostly used when variables
// reference other declarations and such. This remains to be implemented.
func readTypeObj(o *ast.Object) []string {
	// Decl related. Not yet implemented
	vlog(o, "type obj: ")

  if o.Decl == nil {
    return nil
  }
  switch o.Decl.(type) {
  case *ast.TypeSpec:
    t := o.Decl.(*ast.TypeSpec)
    return []string{t.Name.Name}
  }
	return nil
}

func readValueObj(o *ast.Object) (types, values []string) {
	vlog(o, "value obj: ")

  if o.Decl == nil {
    return
  }
  switch o.Decl.(type) {
  case *ast.ValueSpec:
    val := cur_pkg.lookupValue(o.Decl.(*ast.ValueSpec))
    vlog(val, "val: ")
    for i, n := range val.Names {
      if o.Name == n {
        if i < len(val.Types) {
          types = []string{val.Types[i]}
        }
        if i < len(val.Values) {
          values = []string{val.Values[i]}
        }
        return
      }
    }
  }
  return
}
