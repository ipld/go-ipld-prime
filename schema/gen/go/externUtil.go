package gengo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	path "path/filepath"
	"strings"
)

// getExternTypes provides a mapping of all types defined in the destination package.
// It is used by generate to not duplicate defined types to allow overriding of types.
func getExternTypes(pth string) (map[string]struct{}, error) {
	entries, err := os.ReadDir(pth)
	if err != nil {
		return nil, err
	}

	set := token.NewFileSet()
	types := make(map[string]struct{})
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasPrefix(name, "ipldsch_") {
			continue
		}
		f, err := parser.ParseFile(set, path.Join(pth, name), nil, 0)
		if err != nil {
			return nil, err
		}
		for _, d := range f.Decls {
			if t, isType := d.(*ast.GenDecl); isType {
				if t.Tok == token.TYPE {
					for _, s := range t.Specs {
						ts := s.(*ast.TypeSpec)
						types[ts.Name.Name] = struct{}{}
					}
				}
			}
		}
	}

	return types, nil
}
