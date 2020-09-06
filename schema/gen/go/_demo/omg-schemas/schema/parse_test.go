package schema

import (
	"bytes"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
)

// Note! This test file will not *compile*, much less *run*, unless you've invoked codegen.
//  (We're inside a package with "_" prefixed to the name for a reason!
//   `go test ./...` should not generally have gotten you here unawares.)
//  Codegen outputs are not currently committed because we're rapidly iterating on them.

// TestSchemaSchemaParse takes the schema-schema.json document -- the self-describing schema --
// and attempts to unmarshal it into our code-generated schema DMT types.
//
// This is *not* exactly the schema-schema that's upstream in the specs repo -- yet.
// We've made some alterations to make it conform to learnings had during implementing this.
// Some of these alterations may make it back it up to the schema-schema in the specs repo
// (after, of course, sustaining further discussion).
// In particular:
//
// 	- 'TypeDefn' is *keyed* here, whereas it used *inline* in the schema-schema.
//  - a 'Unit' type is introduced (and might belong in the prelude!).
//  - enums are specified using the 'Unit' type (which means serially, they have `{}` instead of `null`).
//
// There's also a couple errata in the programmatic definitions of the types we used for codegen
// where further thought generated an idea which may generate schema-schema changes,
// but these are mostly naming and organizational, and so are minor.
//
func TestSchemaSchemaParse(t *testing.T) {
	nb := Type.Schema__Repr.NewBuilder()
	if err := dagjson.Decoder(nb, bytes.NewBufferString(`
{
	"types": {
		"TypeName": {
			"string": {}
		},
		"SchemaMap": {
			"map": {
				"keyType": "TypeName",
				"valueType": "TypeDefn"
			}
		},
		"AdvancedDataLayoutName": {
			"string": {}
		},
		"AdvancedDataLayoutMap": {
			"map": {
				"keyType": "AdvancedDataLayoutName",
				"valueType": "AdvancedDataLayout"
			}
		},
		"Schema": {
			"struct": {
				"fields": {
					"types": {
						"type": "SchemaMap"
					},
					"advanced": {
						"type": "AdvancedDataLayoutMap"
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"TypeDefn": {
			"union": {
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
						"type": "BytesRepresentation"
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"BytesRepresentation": {
			"union": {
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
						"type": "TypeName"
					},
					"valueType": {
						"type": "TypeNameOrInlineDefn"
					},
					"valueNullable": {
						"type": "Bool"
					},
					"representation": {
						"type": "MapRepresentation"
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
						"type": "String"
					},
					"entryDelim": {
						"type": "String"
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
						"type": "TypeNameOrInlineDefn"
					},
					"valueNullable": {
						"type": "Bool"
					},
					"representation": {
						"type": "ListRepresentation"
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
						"type": "String"
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
						"type": "UnionRepresentation"
					}
				},
				"representation": {
					"map": {}
				}
			}
		},
		"UnionRepresentation": {
			"union": {
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
				"valueType": "TypeName"
			}
		},
		"UnionRepresentation_Keyed": {
			"map": {
				"keyType": "String",
				"valueType": "TypeName"
			}
		},
		"UnionRepresentation_Envelope": {
			"struct": {
				"fields": {
					"discriminantKey": {
						"type": "String"
					},
					"contentKey": {
						"type": "String"
					},
					"discriminantTable": {
						"type": {
							"map": {
								"keyType": "String",
								"valueType": "TypeName"
							}
						}
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
						"type": "String"
					},
					"discriminantTable": {
						"type": {
							"map": {
								"keyType": "String",
								"valueType": "TypeName"
							}
						}
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
								"valueType": "Int"
							}
						}
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
								"valueType": "StructField"
							}
						}
					},
					"representation": {
						"type": "StructRepresentation"
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
						"type": "TypeNameOrInlineDefn"
					},
					"optional": {
						"type": "Bool"
					},
					"nullable": {
						"type": "Bool"
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
								"valueType": "StructRepresentation_Map_FieldDetails"
							}
						},
						"optional": true
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
						"optional": true
					},
					"implicit": {
						"type": "AnyScalar",
						"optional": true
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
								"valueType": "FieldName"
							}
						},
						"optional": true
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
						"type": "String"
					},
					"entryDelim": {
						"type": "String"
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
						"type": "String"
					},
					"fieldOrder": {
						"type": {
							"list": {
								"valueType": "FieldName"
							}
						},
						"optional": true
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
								"valueType": "Unit"
							}
						}
					},
					"representation": {
						"type": "EnumRepresentation"
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
				"valueType": "String"
			}
		},
		"EnumRepresentation_Int": {
			"map": {
				"keyType": "EnumValue",
				"valueType": "Int"
			}
		},
		"TypeCopy": {
			"struct": {
				"fields": {
					"fromType": {
						"type": "TypeName"
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
		panic(err)
	}
}
