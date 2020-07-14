package gengo

import (
	"fmt"
	"strings"

	"github.com/ipld/go-ipld-prime/schema"
)

// This entire file is placeholder-quality implementations.
//
// The AdjunctCfg struct should be replaced with an IPLD Schema-specified thing!
// The values in the unionMemlayout field should be an enum;
// etcetera!

type FieldTuple struct {
	TypeName  schema.TypeName
	FieldName string
}

type AdjunctCfg struct {
	typeSymbolOverrides       map[schema.TypeName]string
	FieldSymbolLowerOverrides map[FieldTuple]string
	fieldSymbolUpperOverrides map[FieldTuple]string
	maybeUsesPtr              map[schema.TypeName]bool   // treat absent as true
	CfgUnionMemlayout         map[schema.TypeName]string // "embedAll"|"interface"; maybe more options later, unclear for now.

	// ... some of these fields have sprouted messy name prefixes so they don't collide with their matching method names.
	//  this structure has reached the critical threshhold where it due to be cleaned up and taken seriously.

	// note: PkgName doesn't appear in here, because it's...
	//  not adjunct data.  it's a generation invocation parameter.
	//   ... this might not hold up in the future though.
	//    There are unanswered questions about how (also, tbf, *if*) we'll handle generation of multiple packages which use each other's types.
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
	if x, ok := cfg.FieldSymbolLowerOverrides[FieldTuple{f.Type().Name(), f.Name()}]; ok {
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

// UnionMemlayout returns a plain string at present;
// there's a case-switch in the templates that processes it.
// We validate that it's a known string when this method is called.
// This should probably be improved in type-safety,
// and validated more aggressively up front when adjcfg is loaded.
func (cfg *AdjunctCfg) UnionMemlayout(t schema.Type) string {
	if t.Kind() != schema.Kind_Union {
		panic(fmt.Errorf("%s is not a union!", t.Name()))
	}
	v, ok := cfg.CfgUnionMemlayout[t.Name()]
	if !ok {
		return "embedAll"
	}
	switch v {
	case "embedAll", "interface":
		return v
	default:
		panic(fmt.Errorf("invalid config: unionMemlayout values must be either \"embedAll\" or \"interface\", not %q", v))
	}
}
