package schemadmt

//go:generate go test -run=Generate -vet=off -tags=bindnodegen

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

// This schema follows https://ipld.io/specs/schemas/schema-schema.ipldsch.

var Type struct {
	Schema schema.TypedPrototype
}

var schemaTypeSystem schema.TypeSystem

func init() {
	var ts schema.TypeSystem
	ts.Init()
	// adjCfg := &gengo.AdjunctCfg{
	// 	FieldSymbolLowerOverrides: map[gengo.FieldTuple]string{
	// 		{"StructField", "type"}: "typ",
	// 	},
	// 	CfgUnionMemlayout: map[schema.TypeName]string{
	// 		"InlineDefn": "interface", // breaks cycles in embeddery that would otherwise be problematic.
	// 	},
	// }

	// I've elided all references to Advancedlayouts stuff for the moment.
	// (Not because it's particularly hard or problematic; I just want to draw a slightly smaller circle first.)

	// Prelude
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnBool("Bool"))
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnFloat("Float"))
	ts.Accumulate(schema.SpawnBytes("Bytes"))

	// Schema-schema!
	// In the same order as the spec's ipldsch file.
	// Note that ADL stuff is excluded for now, as per above.
	ts.Accumulate(schema.SpawnString("TypeName"))
	ts.Accumulate(schema.SpawnMap("SchemaMap",
		"TypeName", "TypeDefn", false,
	))
	ts.Accumulate(schema.SpawnStruct("Schema",
		[]schema.StructField{
			schema.SpawnStructField("types", "SchemaMap", false, false),
			// also: `advanced AdvancedDataLayoutMap`, but as commented above, we'll pursue this later.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("TypeDefn", // TODO: spec's name is "Type"; conflicts with codegen's "var Type typeSlab"
		[]schema.TypeName{
			"TypeBool",
			"TypeString",
			"TypeBytes",
			"TypeInt",
			"TypeFloat",
			"TypeMap",
			"TypeList",
			"TypeLink",
			"TypeUnion",
			"TypeStruct",
			"TypeEnum",
			"TypeCopy",
		},
		// TODO: spec uses inline repr.
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"bool":   "TypeBool",
			"string": "TypeString",
			"bytes":  "TypeBytes",
			"int":    "TypeInt",
			"float":  "TypeFloat",
			"map":    "TypeMap",
			"list":   "TypeList",
			"link":   "TypeLink",
			"union":  "TypeUnion",
			"struct": "TypeStruct",
			"enum":   "TypeEnum",
			"copy":   "TypeCopy",
		}),
	))
	ts.Accumulate(schema.SpawnUnion("TypeTerm",
		[]schema.TypeName{
			"TypeName",
			"InlineDefn",
		},
		schema.SpawnUnionRepresentationKinded(map[datamodel.Kind]schema.TypeName{
			datamodel.Kind_String: "TypeName",
			datamodel.Kind_Map:    "InlineDefn",
		}),
	))
	ts.Accumulate(schema.SpawnUnion("InlineDefn",
		[]schema.TypeName{
			"TypeMap",
			"TypeList",
		},
		// TODO: spec uses inline repr.
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"map":  "TypeMap",
			"list": "TypeList",
		}),
	))
	ts.Accumulate(schema.SpawnStruct("TypeBool",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeString",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeBytes",
		[]schema.StructField{},
		// No BytesRepresentation, since we omit ADL stuff.
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeInt",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeFloat",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeMap",
		[]schema.StructField{
			schema.SpawnStructField("keyType", "TypeName", false, false),
			schema.SpawnStructField("valueType", "TypeTerm", false, false),
			schema.SpawnStructField("valueNullable", "Bool", false, false), // TODO: wants to use the "implicit" feature, but not supported yet
			schema.SpawnStructField("representation", "MapRepresentation", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("MapRepresentation",
		[]schema.TypeName{
			"MapRepresentation_Map",
			"MapRepresentation_Stringpairs",
			"MapRepresentation_Listpairs",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"map":         "MapRepresentation_Map",
			"stringpairs": "MapRepresentation_Stringpairs",
			"listpairs":   "MapRepresentation_Listpairs",
		}),
	))
	ts.Accumulate(schema.SpawnStruct("MapRepresentation_Map",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("MapRepresentation_Stringpairs",
		[]schema.StructField{
			schema.SpawnStructField("innerDelim", "String", false, false),
			schema.SpawnStructField("entryDelim", "String", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("MapRepresentation_Listpairs",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeList",
		[]schema.StructField{
			schema.SpawnStructField("valueType", "TypeTerm", false, false),
			schema.SpawnStructField("valueNullable", "Bool", false, false), // TODO: wants to use the "implicit" feature, but not supported yet
			schema.SpawnStructField("representation", "ListRepresentation", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("ListRepresentation",
		[]schema.TypeName{
			"ListRepresentation_List",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"list": "ListRepresentation_List",
		}),
	))
	ts.Accumulate(schema.SpawnStruct("ListRepresentation_List",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeUnion",
		[]schema.StructField{
			// n.b. we could conceivably allow TypeTerm here rather than just TypeName.  but... we'd rather not: imagine what that means about the type-level behavior of the union: the name munge for the anonymous type would suddenly become load-bearing.  would rather not.
			schema.SpawnStructField("members", "List__TypeName", false, false), // todo: this is a slight hack: should be using an inline defn, but we banged it with name munge coincidents to simplify bootstrap.
			schema.SpawnStructField("representation", "UnionRepresentation", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnList("List__TypeName", // todo: this is a slight hack: should be an anon inside TypeUnion.members.
		"TypeName", false,
	))
	ts.Accumulate(schema.SpawnStruct("TypeLink",
		[]schema.StructField{
			schema.SpawnStructField("expectedType", "TypeName", true, false), // todo: this uses an implicit with a value of 'any' in the schema-schema, but that's been questioned before.  maybe it should simply be an optional.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("UnionRepresentation",
		[]schema.TypeName{
			"UnionRepresentation_Kinded",
			"UnionRepresentation_Keyed",
			"UnionRepresentation_Envelope",
			"UnionRepresentation_Inline",
			"UnionRepresentation_StringPrefix",
			"UnionRepresentation_BytePrefix",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"kinded":       "UnionRepresentation_Kinded",
			"keyed":        "UnionRepresentation_Keyed",
			"envelope":     "UnionRepresentation_Envelope",
			"inline":       "UnionRepresentation_Inline",
			"stringprefix": "UnionRepresentation_StringPrefix",
			"byteprefix":   "UnionRepresentation_BytePrefix",
		}),
	))
	ts.Accumulate(schema.SpawnMap("UnionRepresentation_Kinded",
		"RepresentationKind", "TypeName", false,
	))
	ts.Accumulate(schema.SpawnMap("UnionRepresentation_Keyed",
		"String", "TypeName", false,
	))
	ts.Accumulate(schema.SpawnStruct("UnionRepresentation_Envelope",
		[]schema.StructField{
			schema.SpawnStructField("discriminantKey", "String", false, false),
			schema.SpawnStructField("contentKey", "String", false, false),
			schema.SpawnStructField("discriminantTable", "Map__String__TypeName", false, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("UnionRepresentation_Inline",
		[]schema.StructField{
			schema.SpawnStructField("discriminantKey", "String", false, false),
			schema.SpawnStructField("discriminantTable", "Map__String__TypeName", false, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("UnionRepresentation_StringPrefix",
		[]schema.StructField{
			schema.SpawnStructField("discriminantTable", "Map__String__TypeName", false, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("UnionRepresentation_BytePrefix",
		[]schema.StructField{
			// REVIEW: for schema-schema overall: this is a very funny type.  Should we use strings here?  And perhaps make it use hex for maximum clarity?  This would also allow multi-byte prefixes, which would match what's already done by stringprefix representation.
			schema.SpawnStructField("discriminantTable", "Map__TypeName__Int", false, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnMap("Map__String__TypeName",
		"String", "TypeName", false,
	))
	ts.Accumulate(schema.SpawnMap("Map__TypeName__Int",
		"String", "Int", false,
	))
	ts.Accumulate(schema.SpawnString("RepresentationKind")) // todo: RepresentationKind is supposed to be an enum, but we're puting it to a string atm.
	ts.Accumulate(schema.SpawnStruct("TypeStruct",
		[]schema.StructField{
			schema.SpawnStructField("fields", "Map__FieldName__StructField", false, false), // todo: dodging inline defn's again.
			schema.SpawnStructField("representation", "StructRepresentation", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnMap("Map__FieldName__StructField",
		"FieldName", "StructField", false,
	))
	ts.Accumulate(schema.SpawnString("FieldName"))
	ts.Accumulate(schema.SpawnStruct("StructField",
		[]schema.StructField{
			schema.SpawnStructField("type", "TypeTerm", false, false),
			schema.SpawnStructField("optional", "Bool", false, false), // todo: wants to use the "implicit" feature, but not supported yet
			schema.SpawnStructField("nullable", "Bool", false, false), // todo: wants to use the "implicit" feature, but not supported yet
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("StructRepresentation",
		[]schema.TypeName{
			"StructRepresentation_Map",
			"StructRepresentation_Tuple",
			"StructRepresentation_Stringpairs",
			"StructRepresentation_Stringjoin",
			"StructRepresentation_Listpairs",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"map":         "StructRepresentation_Map",
			"tuple":       "StructRepresentation_Tuple",
			"stringpairs": "StructRepresentation_Stringpairs",
			"stringjoin":  "StructRepresentation_Stringjoin",
			"listpairs":   "StructRepresentation_Listpairs",
		}),
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_Map",
		[]schema.StructField{
			schema.SpawnStructField("fields", "Map__FieldName__StructRepresentation_Map_FieldDetails", true, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnMap("Map__FieldName__StructRepresentation_Map_FieldDetails",
		"FieldName", "StructRepresentation_Map_FieldDetails", false,
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_Map_FieldDetails",
		[]schema.StructField{
			schema.SpawnStructField("rename", "String", true, false),
			schema.SpawnStructField("implicit", "AnyScalar", true, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_Tuple",
		[]schema.StructField{
			schema.SpawnStructField("fieldOrder", "List__FieldName", true, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnList("List__FieldName",
		"FieldName", false,
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_Stringpairs",
		[]schema.StructField{
			schema.SpawnStructField("innerDelim", "String", false, false),
			schema.SpawnStructField("entryDelim", "String", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_Stringjoin",
		[]schema.StructField{
			schema.SpawnStructField("join", "String", false, false),               // review: "delim" would seem more consistent with others -- but this is currently what the schema-schema says.
			schema.SpawnStructField("fieldOrder", "List__FieldName", true, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_Listpairs",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeEnum",
		[]schema.StructField{
			schema.SpawnStructField("members", "Map__EnumValue__Unit", false, false), // todo: dodging inline defn's again.  also: this says unit; schema-schema does not.  schema-schema needs revisiting on this subject.
			schema.SpawnStructField("representation", "EnumRepresentation", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnMap("Map__EnumValue__Unit",
		"EnumValue", "Unit", false,
	))
	ts.Accumulate(schema.SpawnStruct("Unit", // todo: we should formalize the introdution of unit as first class type kind.
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("EnumRepresentation",
		[]schema.TypeName{
			"EnumRepresentation_String",
			"EnumRepresentation_Int",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"string": "EnumRepresentation_String",
			"int":    "EnumRepresentation_Int",
		}),
	))
	ts.Accumulate(schema.SpawnString("EnumValue"))
	ts.Accumulate(schema.SpawnMap("EnumRepresentation_String",
		"EnumValue", "String", false,
	))
	ts.Accumulate(schema.SpawnMap("EnumRepresentation_Int",
		"EnumValue", "Int", false,
	))
	ts.Accumulate(schema.SpawnStruct("TypeCopy",
		[]schema.StructField{
			schema.SpawnStructField("fromType", "TypeName", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("AnyScalar",
		[]schema.TypeName{
			"Bool",
			"String",
			"Bytes",
			"Int",
			"Float",
		},
		schema.SpawnUnionRepresentationKinded(map[datamodel.Kind]schema.TypeName{
			datamodel.Kind_Bool:   "Bool",
			datamodel.Kind_String: "String",
			datamodel.Kind_Bytes:  "Bytes",
			datamodel.Kind_Int:    "Int",
			datamodel.Kind_Float:  "Float",
		}),
	))

	if errs := ts.ValidateGraph(); errs != nil {
		for _, err := range errs {
			fmt.Printf("- %s\n", err)
		}
		panic("not happening")
	}

	schemaTypeSystem = ts

	Type.Schema = bindnode.Prototype(
		(*Schema)(nil),
		schemaTypeSystem.TypeByName("Schema"),
	)
}
