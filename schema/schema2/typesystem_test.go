package schema

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

func TestBuildTypeSystem(t *testing.T) {
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
		Wish(t, ts.types["Woop"], ShouldBeSameTypeAs, &TypeString{})
		Wish(t, ts.types["Woop"].TypeKind(), ShouldEqual, TypeKind_String)
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
				fmt.Errorf("type SomeList refers to missing type Bork as value type"),
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
				fmt.Errorf("type SomeMap refers to missing type Bork as key type"),
				fmt.Errorf("type SomeMap refers to missing type Spork as value type"),
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
		Wish(t, ts.types["SomeMap"], ShouldBeSameTypeAs, &TypeMap{})
		Wish(t, ts.types["SomeMap"].TypeKind(), ShouldEqual, TypeKind_Map)
		Wish(t, ts.types["SomeMap"].(*TypeMap).KeyType().Name().String(), ShouldEqual, "String")
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
		Wish(t, ts.types["SomeMap"], ShouldBeSameTypeAs, &TypeMap{})
		Wish(t, ts.types["SomeMap"].TypeKind(), ShouldEqual, TypeKind_Map)
		Wish(t, ts.types["SomeMap"].(*TypeMap).KeyType().Name().String(), ShouldEqual, "StringyStruct")
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

func testParse(t *testing.T, schemajson string, expectParseErr error, expectTypesystemError []error) *TypeSystem {
	t.Helper()
	dmt, parseErr := parseSchema(schemajson)
	Wish(t, parseErr, ShouldEqual, expectParseErr)
	if parseErr != nil {
		return nil
	}
	ts, typesystemErr := BuildTypeSystem(dmt)
	Wish(t, typesystemErr, ShouldEqual, expectTypesystemError)
	return ts
}

func parseSchema(schemajson string) (schemadmt.Schema, error) {
	nb := schemadmt.Type.Schema__Repr.NewBuilder()
	if err := dagjson.Decode(nb, strings.NewReader(schemajson)); err != nil {
		return nil, err
	}
	return nb.Build().(schemadmt.Schema), nil
}
