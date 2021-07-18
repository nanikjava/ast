package main

import (
	"bytes"
	"fmt"
	"github.com/karrick/godirwalk"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func fix(dir string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			fmt.Printf("working on file %v\n", fileName)
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if ok {
					fmt.Printf("fn: %v\n", fn.Name.Name)

					if fn.Name.Obj != nil {
						fDeclaration, ok := fn.Name.Obj.Decl.(*ast.FuncDecl)
						if ok {
							if len(fDeclaration.Type.Params.List) > 0 {
								// how many parameters are there
								fmt.Println("Len param  = ", len(fDeclaration.Type.Params.List))

								for _, k := range fDeclaration.Type.Params.List {
									if len(k.Names) > 0 {
										fmt.Println(". parameter name = ", k.Names[0])

										// is it function type ?
										fType, ok := k.Type.(*ast.FuncType)
										if fType != nil && ok {
											fmt.Println(".... parameter type (function) = ", fType.Params.List, " with ", len(fType.Params.List), " parameters ")
										}

										// is it a selector type ?
										fSelectorExpr, ok := k.Type.(*ast.SelectorExpr)
										if fSelectorExpr != nil && ok {
											fIdent, ok := fSelectorExpr.X.(*ast.Ident)
											if ok {
												fmt.Println(".... parameter type = ", fIdent.Name)
											}
										}

										// ..or just a nomal type ?
										fIdent, ok := k.Type.(*ast.Ident)
										if ok {
											if fIdent != nil && fIdent.Obj != nil && ok {
												fmt.Println(".... parameter type = ", fIdent.Obj.Name)
											} else {
												fmt.Println(".... parameter type = ", fIdent.Name)
											}
										}
									}
								}
							}
						}
					}
				}
				return true
			})

			buf := new(bytes.Buffer)
			err := format.Node(buf, fset, file)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			} else if fileName[len(fileName)-8:] != "_test.go" {
				ioutil.WriteFile(fileName, buf.Bytes(), 0664)
			}
		}
	}
}

func main() {
	traverse("/home/nanik/GolandProjects/gopath/src/github.com/ory/kratos")
}

func traverse(mainDir string) {
	godirwalk.Walk(mainDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			// Following string operation is not most performant way
			// of doing this, but common enough to warrant a simple
			// example here:
			if de.IsDir() {
				fmt.Printf("%s %s\n", de.ModeType(), osPathname)
				fix(osPathname)
			}
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
}
