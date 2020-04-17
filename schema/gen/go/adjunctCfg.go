package gengo

import (
	"strings"

	"github.com/ipld/go-ipld-prime/schema"
)

type FieldTuple struct {
	TypeName  schema.TypeName
	FieldName string
}

type AdjunctCfg struct {
	typeSymbolOverrides       map[schema.TypeName]string
	fieldSymbolLowerOverrides map[FieldTuple]string
	fieldSymbolUpperOverrides map[FieldTuple]string
	maybeUsesPtr              map[schema.TypeName]bool // treat absent as true

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
	if x, ok := cfg.fieldSymbolLowerOverrides[FieldTuple{f.Type().Name(), f.Name()}]; ok {
		return x
	}
	return f.Name() // presumed already lower
}

func (cfg *AdjunctCfg) FieldSymbolUpper(f schema.StructField) string {
	if x, ok := cfg.fieldSymbolUpperOverrides[FieldTuple{f.Type().Name(), f.Name()}]; ok {
		return x
	}
	return strings.Title(f.Name())
}

func (cfg *AdjunctCfg) MaybeUsesPtr(t schema.Type) bool {
	if x, ok := cfg.maybeUsesPtr[t.Name()]; ok {
		return x
	}
	// FUTURE: we could make this default vary based on sizeof the type.
	//  It's generally true that for scalars it should be false by default; and that's easy to do.
	//   It would actually *remove* special cases from the prelude, which would be a win.
	//  Maps and lists should also probably default off...?
	//   (I have a feeling something might get touchy there.  Review when implementing those.)
	//  Perhaps structs and unions are the only things likely to benefit from pointers.
	return true
}
