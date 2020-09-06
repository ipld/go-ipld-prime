package schemaschema

import (
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
	gengo "github.com/ipld/go-ipld-prime/schema/gen/go"
)

func init() {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &gengo.AdjunctCfg{
		FieldSymbolLowerOverrides: map[gengo.FieldTuple]string{
			{"StructField", "type"}: "typ",
			{"TypeEnum", "type"}:    "typ",
		},
		CfgUnionMemlayout: map[schema.TypeName]string{
			"TypeDefnInline": "interface", // breaks cycles in embeddery that would otherwise be problematic.
		},
	}

	// I've elided all references to Advancedlayouts stuff for the moment.
	// (Not because it's particularly hard or problematic; I just want to draw a slightly smaller circle first.)

	// Prelude
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnBool("Bool"))
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnFloat("Float"))
	ts.Accumulate(schema.SpawnBytes("Bytes"))

	// Schema-schema!
	ts.Accumulate(schema.SpawnStruct("Schema",
		[]schema.StructField{
			schema.SpawnStructField("types", "SchemaMap", false, false),
			// also: `advanced AdvancedDataLayoutMap`, but as commented above, we'll pursue this later.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnString("TypeName"))
	ts.Accumulate(schema.SpawnMap("SchemaMap",
		"TypeName", "TypeDefn", false,
	))
	ts.Accumulate(schema.SpawnUnion("TypeDefn",
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
	ts.Accumulate(schema.SpawnUnion("TypeNameOrInlineDefn",
		[]schema.TypeName{
			"TypeName",
			"TypeDefnInline",
		},
		schema.SpawnUnionRepresentationKinded(map[ipld.ReprKind]schema.TypeName{
			ipld.ReprKind_String: "TypeName",
			ipld.ReprKind_Map:    "TypeDefnInline",
		}),
	))
	ts.Accumulate(schema.SpawnUnion("TypeDefnInline", // n.b. previously called "TypeTerm".
		[]schema.TypeName{
			"TypeMap",
			"TypeList",
		},
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
	ts.Accumulate(schema.SpawnStruct("TypeLink",
		[]schema.StructField{
			schema.SpawnStructField("expectedType", "TypeName", true, false), // todo: this uses an implicit with a value of 'any' in the schema-schema, but that's been questioned before.  maybe it should simply be an optional.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeMap",
		[]schema.StructField{
			schema.SpawnStructField("keyType", "TypeName", false, false),
			schema.SpawnStructField("valueType", "TypeNameOrInlineDefn", false, false),
			schema.SpawnStructField("valueNullable", "Bool", false, false), // todo: wants to use the "implicit" feature, but not supported yet
			schema.SpawnStructField("representation", "MapRepresentation", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("MapRepresentation",
		[]schema.TypeName{
			"MapRepresentation_Map",
			"MapRepresentation_StringPairs",
			"MapRepresentation_ListPairs",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"map":         "MapRepresentation_Map",
			"stringpairs": "MapRepresentation_StringPairs",
			"listpairs":   "MapRepresentation_ListPairs",
		}),
	))
	ts.Accumulate(schema.SpawnStruct("MapRepresentation_Map",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("MapRepresentation_StringPairs",
		[]schema.StructField{
			schema.SpawnStructField("innerDelim", "String", false, false),
			schema.SpawnStructField("entryDelim", "String", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("MapRepresentation_ListPairs",
		[]schema.StructField{},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("TypeList",
		[]schema.StructField{
			schema.SpawnStructField("valueType", "TypeNameOrInlineDefn", false, false),
			schema.SpawnStructField("valueNullable", "Bool", false, false), // todo: wants to use the "implicit" feature, but not supported yet
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
			// n.b. we could conceivably allow TypeNameOrInlineDefn here rather than just TypeName.  but... we'd rather not: imagine what that means about the type-level behavior of the union: the name munge for the anonymous type would suddenly become load-bearing.  would rather not.
			schema.SpawnStructField("members", "List__TypeName", false, false), // todo: this is a slight hack: should be using an inline defn, but we banged it with name munge coincidents to simplify bootstrap.
			schema.SpawnStructField("representation", "UnionRepresentation", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnList("List__TypeName", // todo: this is a slight hack: should be an anon inside TypeUnion.members.
		"TypeName", false,
	))
	ts.Accumulate(schema.SpawnUnion("UnionRepresentation",
		[]schema.TypeName{
			"UnionRepresentation_Kinded",
			"UnionRepresentation_Keyed",
			"UnionRepresentation_Envelope",
			"UnionRepresentation_Inline",
			"UnionRepresentation_BytePrefix",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"kinded":     "UnionRepresentation_Kinded",
			"keyed":      "UnionRepresentation_Keyed",
			"envelope":   "UnionRepresentation_Envelope",
			"inline":     "UnionRepresentation_Inline",
			"byteprefix": "UnionRepresentation_BytePrefix",
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
	ts.Accumulate(schema.SpawnStruct("UnionRepresentation_BytePrefix",
		[]schema.StructField{
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
			schema.SpawnStructField("type", "TypeNameOrInlineDefn", false, false),
			schema.SpawnStructField("optional", "Bool", false, false), // todo: wants to use the "implicit" feature, but not supported yet
			schema.SpawnStructField("nullable", "Bool", false, false), // todo: wants to use the "implicit" feature, but not supported yet
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnUnion("StructRepresentation",
		[]schema.TypeName{
			"StructRepresentation_Map",
			"StructRepresentation_Tuple",
			"StructRepresentation_StringPairs",
			"StructRepresentation_StringJoin",
			"StructRepresentation_ListPairs",
		},
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"map":         "StructRepresentation_Map",
			"tuple":       "StructRepresentation_Tuple",
			"stringpairs": "StructRepresentation_StringPairs",
			"stringjoin":  "StructRepresentation_StringJoin",
			"listpairs":   "StructRepresentation_ListPairs",
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
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_StringPairs",
		[]schema.StructField{
			schema.SpawnStructField("innerDelim", "String", false, false),
			schema.SpawnStructField("entryDelim", "String", false, false),
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_StringJoin",
		[]schema.StructField{
			schema.SpawnStructField("join", "String", false, false),               // review: "delim" would seem more consistent with others -- but this is currently what the schema-schema says.
			schema.SpawnStructField("fieldOrder", "List__FieldName", true, false), // todo: dodging inline defn's again.
		},
		schema.StructRepresentation_Map{},
	))
	ts.Accumulate(schema.SpawnStruct("StructRepresentation_ListPairs",
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
		schema.SpawnUnionRepresentationKinded(map[ipld.ReprKind]schema.TypeName{
			ipld.ReprKind_Bool:   "Bool",
			ipld.ReprKind_String: "String",
			ipld.ReprKind_Bytes:  "Bytes",
			ipld.ReprKind_Int:    "Int",
			ipld.ReprKind_Float:  "Float",
		}),
	))

	if errs := ts.ValidateGraph(); errs != nil {
		for _, err := range errs {
			fmt.Printf("- %s\n", err)
		}
		panic("not happening")
	}

	os.Mkdir("./schema", 0755)
	gengo.Generate("./schema", "schema", ts, adjCfg)
}
