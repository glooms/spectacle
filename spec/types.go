package spec

import (
  "go/ast"
  "go/importer"
  "go/parser"
  "go/token"
  "go/types"
)


func typeCheck(fset *token.FileSet, astPkg *ast.Package) (map[ast.Expr]types.TypeAndValue, error) {
  path := astPkg.Name
  conf := &types.Config{Importer: importer.Default()}
  files := getFiles(astPkg)
  info := &types.Info{
    Types: map[ast.Expr]types.TypeAndValue{},
  }
  pkg, err := conf.Check(path, fset, files, info)
  if err != nil {
    return nil, err
  }
  checker := types.NewChecker(conf, fset, pkg, info)
  _ = checker.Files(files)
  /*
  for k, v := range info.Types {
    s := v.Type.String()
    if v.Value != nil {
      s += ": " + v.Value.String()
    }
    vlog(k)
    log(s)
  } */
  return info.Types, nil
}

func getFiles(pkg *ast.Package) []*ast.File {
  files := make([]*ast.File, len(pkg.Files))
  i := 0
  for _, file := range pkg.Files {
    files[i] = file
    i++
  }
  return files
}

func test() {
  path := "decls"
  fset := token.NewFileSet()
  pkgs, err := parser.ParseDir(fset, path, filter, 0)
  check(err)
  typeCheck(fset, pkgs["decls"])
}
