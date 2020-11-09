package gengraphqlserver

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
)

func GetPreExistingMethods(dst string) map[string]struct{} {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, dst, nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package for overrides: ", err)
		return make(map[string]struct{})
	}

	funcNames := make(map[string]struct{})
	for _, pack := range packs {
		for fname, f := range pack.Files {
			if path.Base(fname) == "schema.go" {
				continue
			}
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					funcNames[fn.Name.Name] = struct{}{}
				}
			}
		}
	}

	return funcNames
}
