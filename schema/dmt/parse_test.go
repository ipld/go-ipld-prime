package schemadmt_test

import (
	"strings"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	// "github.com/ipld/go-ipld-prime/schema/schema2"
)

// TestSchemaSchemaParse takes the schema-schema.json document -- the self-describing schema --
// and attempts to unmarshal it into our code-generated schema DMT types.
//
// This is *not* exactly the schema-schema that's upstream in the specs repo -- yet.
// We've made some alterations to make it conform to learnings had during implementing this.
// Some of these alterations may make it back it up to the schema-schema in the specs repo
// (after, of course, sustaining further discussion).
//
// The changes that might be worth upstreaming are:
//
// 	- 'TypeDefn' is *keyed* here, whereas it used *inline* in the schema-schema.
//  - a 'Unit' type is introduced (and might belong in the prelude!).
//  - enums are specified using the 'Unit' type (which means serially, they have `{}` instead of `null`).
//  - a few naming changes, which are minor and nonsemantic.
//
// There's also a few accumulated changes which are working around incomplete
// features of our own tooling here, and are bugs that should be fixed (definitely not upstreamed):
//
//  - many field definitions have a `"optional": false, "nullable": false`
//    explicitly stated, where it should be sufficient to leave these implicit.
//    (These are avoiding our current lack of support for implicits.)
//  - similarly, many map definitions have an `"valueNullable": false`
//    explicitly stated, where it should be sufficient to leave these implicit.
//
func TestSchemaSchemaParse(t *testing.T) {
	nb := schemadmt.Type.Schema__Repr.NewBuilder()
	if err := dagjson.Decode(nb, strings.NewReader(`
{
	"types": {
		"TypeName": {
			"string": {}
		},
		"SchemaMap": {
			"map": {
				"keyType": "TypeName",
				"valueType": "TypeDefn",
				"valueNullable": false,
				"representation": {
					"map":{}
				}
			}
		},
		"AdvancedDataLayoutName": {
			"string": {}
		},
		"AdvancedDataLayoutMap": {
			"map": {
				"keyType": "AdvancedDataLayoutName",
				"valueType": "AdvancedDataLayout",
				"valueNullable": false,
				"representation": {
					"map":{}
				}
			}
		},
		"Schema": {
			"struct": {
				"fields": {
					"types": {
						"type": "SchemaMap",
						"optional": false,
						"nullable": false
					},
					"advanced": {
						"type": "AdvancedDataLayoutMap",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeDefn": {
			"union": {
				"members": [
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
					"TypeCopy"
				],
				"representation": {
					"keyed": {
						"bool": "TypeBool",
						"string": "TypeString",
						"bytes": "TypeBytes",
						"int": "TypeInt",
						"float": "TypeFloat",
						"map": "TypeMap",
						"list": "TypeList",
						"link": "TypeLink",
						"union": "TypeUnion",
						"struct": "TypeStruct",
						"enum": "TypeEnum",
						"copy": "TypeCopy"
					}
				}
			}
		},
		"TypeKind": {
			"enum": {
				"members": {
					"Bool": {},
					"String": {},
					"Bytes": {},
					"Int": {},
					"Float": {},
					"Map": {},
					"List": {},
					"Link": {},
					"Union": {},
					"Struct": {},
					"Enum": {}
				},
				"representation": {
					"string": {}
				}
			}
		},
		"RepresentationKind": {
			"enum": {
				"members": {
					"Bool": {},
					"String": {},
					"Bytes": {},
					"Int": {},
					"Float": {},
					"Map": {},
					"List": {},
					"Link": {}
				},
				"representation": {
					"string": {}
				}
			}
		},
		"AnyScalar": {
			"union": {
				"members": [
					"Bool",
					"String",
					"Bytes",
					"Int",
					"Float"
				],
				"representation": {
					"kinded": {
						"bool": "Bool",
						"string": "String",
						"bytes": "Bytes",
						"int": "Int",
						"float": "Float"
					}
				}
			}
		},
		"AdvancedDataLayout": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeBool": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeString": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeBytes": {
			"struct": {
				"fields": {
					"representation": {
						"type": "BytesRepresentation",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"BytesRepresentation": {
			"union": {
				"members": [
					"BytesRepresentation_Bytes",
					"AdvancedDataLayoutName"
				],
				"representation": {
					"keyed": {
						"bytes": "BytesRepresentation_Bytes",
						"advanced": "AdvancedDataLayoutName"
					}
				}
			}
		},
		"BytesRepresentation_Bytes": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeInt": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeFloat": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeMap": {
			"struct": {
				"fields": {
					"keyType": {
						"type": "TypeName",
						"optional": false,
						"nullable": false
					},
					"valueType": {
						"type": "TypeNameOrInlineDefn",
						"optional": false,
						"nullable": false
					},
					"valueNullable": {
						"type": "Bool",
						"optional": false,
						"nullable": false
					},
					"representation": {
						"type": "MapRepresentation",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {
						"fields": {
							"valueNullable": {
								"implicit": false
							}
						}
					}
				}
			}
		},
		"MapRepresentation": {
			"union": {
				"members": [
					"MapRepresentation_Map",
					"MapRepresentation_StringPairs",
					"MapRepresentation_ListPairs",
					"AdvancedDataLayoutName"
				],
				"representation": {
					"keyed": {
						"map": "MapRepresentation_Map",
						"stringpairs": "MapRepresentation_StringPairs",
						"listpairs": "MapRepresentation_ListPairs",
						"advanced": "AdvancedDataLayoutName"
					}
				}
			}
		},
		"MapRepresentation_Map": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"MapRepresentation_StringPairs": {
			"struct": {
				"fields": {
					"innerDelim": {
						"type": "String",
						"optional": false,
						"nullable": false
					},
					"entryDelim": {
						"type": "String",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"MapRepresentation_ListPairs": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeList": {
			"struct": {
				"fields": {
					"valueType": {
						"type": "TypeNameOrInlineDefn",
						"optional": false,
						"nullable": false
					},
					"valueNullable": {
						"type": "Bool",
						"optional": false,
						"nullable": false
					},
					"representation": {
						"type": "ListRepresentation",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {
						"fields": {
							"valueNullable": {
								"implicit": false
							}
						}
					}
				}
			}
		},
		"ListRepresentation": {
			"union": {
				"members": [
					"ListRepresentation_List",
					"AdvancedDataLayoutName"
				],
				"representation": {
					"keyed": {
						"list": "ListRepresentation_List",
						"advanced": "AdvancedDataLayoutName"
					}
				}
			}
		},
		"ListRepresentation_List": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeLink": {
			"struct": {
				"fields": {
					"expectedType": {
						"type": "String",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {
						"fields": {
							"expectedType": {
								"implicit": "Any"
							}
						}
					}
				}
			}
		},
		"TypeUnion": {
			"struct": {
				"fields": {
					"representation": {
						"type": "UnionRepresentation",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"UnionRepresentation": {
			"union": {
				"members": [
					"UnionRepresentation_Kinded",
					"UnionRepresentation_Keyed",
					"UnionRepresentation_Envelope",
					"UnionRepresentation_Inline",
					"UnionRepresentation_BytePrefix"
				],
				"representation": {
					"keyed": {
						"kinded": "UnionRepresentation_Kinded",
						"keyed": "UnionRepresentation_Keyed",
						"envelope": "UnionRepresentation_Envelope",
						"inline": "UnionRepresentation_Inline",
						"byteprefix": "UnionRepresentation_BytePrefix"
					}
				}
			}
		},
		"UnionRepresentation_Kinded": {
			"map": {
				"keyType": "RepresentationKind",
				"valueType": "TypeName",
				"valueNullable": false,
				"representation": {
					"map":{}
				}
			}
		},
		"UnionRepresentation_Keyed": {
			"map": {
				"keyType": "String",
				"valueType": "TypeName",
				"valueNullable": false,
				"representation": {
					"map":{}
				}
			}
		},
		"UnionRepresentation_Envelope": {
			"struct": {
				"fields": {
					"discriminantKey": {
						"type": "String",
						"optional": false,
						"nullable": false
					},
					"contentKey": {
						"type": "String",
						"optional": false,
						"nullable": false
					},
					"discriminantTable": {
						"type": {
							"map": {
								"keyType": "String",
								"valueType": "TypeName",
								"valueNullable": false,
								"representation": {
									"map":{}
								}
							}
						},
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"UnionRepresentation_Inline": {
			"struct": {
				"fields": {
					"discriminantKey": {
						"type": "String",
						"optional": false,
						"nullable": false
					},
					"discriminantTable": {
						"type": {
							"map": {
								"keyType": "String",
								"valueType": "TypeName",
								"valueNullable": false,
								"representation": {
									"map":{}
								}
							}
						},
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"UnionRepresentation_BytePrefix": {
			"struct": {
				"fields": {
					"discriminantTable": {
						"type": {
							"map": {
								"keyType": "TypeName",
								"valueType": "Int",
								"valueNullable": false,
								"representation": {
									"map":{}
								}
							}
						},
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeStruct": {
			"struct": {
				"fields": {
					"fields": {
						"type": {
							"map": {
								"keyType": "FieldName",
								"valueType": "StructField",
								"valueNullable": false,
								"representation": {
									"map":{}
								}
							}
						},
						"optional": false,
						"nullable": false
					},
					"representation": {
						"type": "StructRepresentation",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"FieldName": {
			"string": {}
		},
		"StructField": {
			"struct": {
				"fields": {
					"type": {
						"type": "TypeNameOrInlineDefn",
						"optional": false,
						"nullable": false
					},
					"optional": {
						"type": "Bool",
						"optional": false,
						"nullable": false
					},
					"nullable": {
						"type": "Bool",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {
						"fields": {
							"optional": {
								"implicit": false
							},
							"nullable": {
								"implicit": false
							}
						}
					}
				}
			}
		},
		"TypeNameOrInlineDefn": {
			"union": {
				"members": [
					"TypeName",
					"TypeDefnInline"
				],
				"representation": {
					"kinded": {
						"string": "TypeName",
						"map": "TypeDefnInline"
					}
				}
			}
		},
		"TypeDefnInline": {
			"union": {
				"members": [
					"TypeMap",
					"TypeList"
				],
				"representation": {
					"keyed": {
						"map": "TypeMap",
						"list": "TypeList"
					}
				}
			}
		},
		"StructRepresentation": {
			"union": {
				"members": [
					"StructRepresentation_Map",
					"StructRepresentation_Tuple",
					"StructRepresentation_StringPairs",
					"StructRepresentation_StringJoin",
					"StructRepresentation_ListPairs"
				],
				"representation": {
					"keyed": {
						"map": "StructRepresentation_Map",
						"tuple": "StructRepresentation_Tuple",
						"stringpairs": "StructRepresentation_StringPairs",
						"stringjoin": "StructRepresentation_StringJoin",
						"listpairs": "StructRepresentation_ListPairs"
					}
				}
			}
		},
		"StructRepresentation_Map": {
			"struct": {
				"fields": {
					"fields": {
						"type": {
							"map": {
								"keyType": "FieldName",
								"valueType": "StructRepresentation_Map_FieldDetails",
								"valueNullable": false,
								"representation": {
									"map":{}
								}
							}
						},
						"optional": true,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"StructRepresentation_Map_FieldDetails": {
			"struct": {
				"fields": {
					"rename": {
						"type": "String",
						"optional": true,
						"nullable": false
					},
					"implicit": {
						"type": "AnyScalar",
						"optional": true,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"StructRepresentation_Tuple": {
			"struct": {
				"fields": {
					"fieldOrder": {
						"type": {
							"list": {
								"valueType": "FieldName",
								"valueNullable": false,
								"representation": {
									"list":{}
								}
							},
						},
						"optional": true,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"StructRepresentation_StringPairs": {
			"struct": {
				"fields": {
					"innerDelim": {
						"type": "String",
						"optional": false,
						"nullable": false
					},
					"entryDelim": {
						"type": "String",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"StructRepresentation_StringJoin": {
			"struct": {
				"fields": {
					"join": {
						"type": "String",
						"optional": false,
						"nullable": false
					},
					"fieldOrder": {
						"type": {
							"list": {
								"valueType": "FieldName",
								"valueNullable": false,
								"representation": {
									"list":{}
								}
							}
						},
						"optional": true,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"StructRepresentation_ListPairs": {
			"struct": {
				"fields": {},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeEnum": {
			"struct": {
				"fields": {
					"members": {
						"type": {
							"map": {
								"keyType": "EnumValue",
								"valueType": "Unit",
								"valueNullable": false,
								"representation": {
									"map":{}
								}
							}
						},
						"optional": false,
						"nullable": false
					},
					"representation": {
						"type": "EnumRepresentation",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"EnumValue": {
			"string": {}
		},
		"EnumRepresentation": {
			"union": {
				"members": [
					"EnumRepresentation_String",
					"EnumRepresentation_Int"
				],
				"representation": {
					"keyed": {
						"string": "EnumRepresentation_String",
						"int": "EnumRepresentation_Int"
					}
				}
			}
		},
		"EnumRepresentation_String": {
			"map": {
				"keyType": "EnumValue",
				"valueType": "String",
				"valueNullable": false,
				"representation": {
					"map":{}
				}
			}
		},
		"EnumRepresentation_Int": {
			"map": {
				"keyType": "EnumValue",
				"valueType": "Int",
				"valueNullable": false,
				"representation": {
					"map":{}
				}
			}
		},
		"TypeCopy": {
			"struct": {
				"fields": {
					"fromType": {
						"type": "TypeName",
						"optional": false,
						"nullable": false
					}
				},
				"representation": {
					"map": {}
				}
			}
		}
	}
}
	`)); err != nil {
		t.Error(err)
	}
	// n := nb.Build().(schemadmt.Schema)

	// Reify that thang!
	// TODO: not yet :) anonymous types used in the above data are not yet implemented.
	// _, errs := schema.BuildTypeSystem(n)
	// if errs != nil {
	// t.Error(errs)
	// }
}
