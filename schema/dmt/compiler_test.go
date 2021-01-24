package schemadmt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/polydawn/refmt/json"
	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestCompile(t *testing.T) {
	// NOTE: several of these fixtures will need updating when support for implicits is completed.
	t.Run("SimpleHappyPath", func(t *testing.T) {
		ts := testParse(t,
			`{
				"types": {
					"Woop": {
						"string": {}
					}
				}
			}`,
			nil,
			nil,
		)
		Wish(t, ts.GetType("Woop"), ShouldBeSameTypeAs, &schema.TypeString{})
		Wish(t, ts.GetType("Woop").TypeKind(), ShouldEqual, schema.TypeKind_String)
	})
	t.Run("MissingTypeInList", func(t *testing.T) {
		testParse(t,
			`{
				"types": {
					"SomeList": {
						"list": {
							"valueType": "Bork",
							"valueNullable": false,
							"representation": {
								"list": {}
							}
						}
					}
				}
			}`,
			nil,
			[]error{
				fmt.Errorf(`type SomeList is invalid: list declaration's value type must be defined: missing type "Bork"`),
			},
		)
	})
	t.Run("MissingTypeInMap", func(t *testing.T) {
		testParse(t,
			`{
				"types": {
					"SomeMap": {
						"map": {
							"keyType": "Bork"
							"valueType": "Spork",
							"valueNullable": false,
							"representation": {
								"map": {}
							}
						}
					}
				}
			}`,
			nil,
			[]error{
				fmt.Errorf(`type SomeMap is invalid: map declaration's key type must be defined: missing type "Bork"`),
				// REVIEW: this is a case where the short-circuit exiting during rule evaluation blocks an easy win:
				//fmt.Errorf(`type SomeMap is invalid: map declaration's value type must be defined: missing type "Spork"`),
			},
		)
	})
	t.Run("SimpleValidMapKeyType", func(t *testing.T) {
		ts := testParse(t,
			`{
				"types": {
					"SomeMap": {
						"map": {
							"keyType": "String"
							"valueType": "String",
							"valueNullable": false,
							"representation": {
								"map": {}
							}
						}
					},
					"String": {
						"string": {}
					}
				}
			}`,
			nil,
			nil,
		)
		Wish(t, ts.GetType("SomeMap"), ShouldBeSameTypeAs, &schema.TypeMap{})
		Wish(t, ts.GetType("SomeMap").TypeKind(), ShouldEqual, schema.TypeKind_Map)
		Wish(t, ts.GetType("SomeMap").(*schema.TypeMap).KeyType().Name().String(), ShouldEqual, "String")
	})
	t.Run("ComplexValidMapKeyType", func(t *testing.T) {
		ts := testParse(t,
			`{
				"types": {
					"SomeMap": {
						"map": {
							"keyType": "StringyStruct",
							"valueType": "String",
							"valueNullable": false,
							"representation": {
								"map": {}
							}
						}
					},
					"String": {
						"string": {}
					},
					"StringyStruct": {
						"struct": {
							"fields": {
								"f1": {
									"type": "String",
									"optional": false,
									"nullable": false
								},
								"f2": {
									"type": "String",
									"optional": false,
									"nullable": false
								}
							},
							"representation": {
								"stringjoin": {
									"join": ":"
								}
							}
						}
					}
				}
			}`,
			nil,
			nil,
		)
		Wish(t, ts.GetType("SomeMap"), ShouldBeSameTypeAs, &schema.TypeMap{})
		Wish(t, ts.GetType("SomeMap").TypeKind(), ShouldEqual, schema.TypeKind_Map)
		Wish(t, ts.GetType("SomeMap").(*schema.TypeMap).KeyType().Name().String(), ShouldEqual, "StringyStruct")
	})
	t.Run("InvalidMapKeyType", func(t *testing.T) {
		testParse(t,
			`{
				"types": {
					"SomeMap": {
						"map": {
							"keyType": "StringyStruct",
							"valueType": "String",
							"valueNullable": false,
							"representation": {
								"map": {}
							}
						}
					},
					"String": {
						"string": {}
					},
					"StringyStruct": {
						"struct": {
							"fields": {
								"f1": {
									"type": "String",
									"optional": false,
									"nullable": false
								},
								"f2": {
									"type": "String",
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
			}`,
			nil,
			[]error{
				fmt.Errorf("type SomeMap refers to type StringyStruct as key type, but it is not a valid key type because it is not stringable"),
			},
		)
	})
}

func testParse(t *testing.T, schemajson string, expectParseErr error, expectTypesystemError []error) *schema.TypeSystem {
	t.Helper()
	dmt, parseErr := parseSchema(schemajson)
	Wish(t, parseErr, ShouldEqual, expectParseErr)
	if parseErr != nil {
		return nil
	}
	ts, typesystemErrs := dmt.Compile()
	Require(t, typesystemErrs, ShouldEqual, expectTypesystemError)
	return ts
}

func parseSchema(schemajson string) (Schema, error) {
	nb := Type.Schema__Repr.NewBuilder()
	if err := dagjson.Unmarshal(nb, json.NewDecoder(strings.NewReader(schemajson))); err != nil {
		return nil, err
	}
	return nb.Build().(Schema), nil
}
