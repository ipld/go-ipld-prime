package gengo

import (
	"strings"

	"github.com/ipld/go-ipld-prime/schema"
)

type AdjunctCfg struct {
	typeSymbolOverrides       map[schema.TypeName]string
	fieldSymbolLowerOverrides map[schema.StructField]string
	fieldSymbolUpperOverrides map[schema.StructField]string

	// note: PkgName doesn't appear in here, because it's...
	//  not adjunct data.  it's a generation invocation parameter.
}

// TypeSymbol returns the symbol for a type;
// by default, it's the same string as its name in the schema,
// but it can be overriden.
//
// This is the base, unembellished symbol.
// It's frequently augmented:
// prefixing an underscore to make it unexported;
// suffixing "__Something" to make the name of a supporting type;
// etc.
// (Most such augmentations are not configurable.)
func (cfg *AdjunctCfg) TypeSymbol(t schema.Type) string {
	if x, ok := cfg.typeSymbolOverrides[t.Name()]; ok {
		return x
	}
	return string(t.Name()) // presumed already upper
}

func (cfg *AdjunctCfg) FieldSymbolLower(f schema.StructField) string {
	if x, ok := cfg.fieldSymbolLowerOverrides[f]; ok {
		return x
	}
	return f.Name() // presumed already lower
}

func (cfg *AdjunctCfg) FieldSymbolUpper(f schema.StructField) string {
	if x, ok := cfg.fieldSymbolUpperOverrides[f]; ok {
		return x
	}
	return strings.Title(f.Name())
}
