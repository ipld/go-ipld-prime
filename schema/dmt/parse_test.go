package schemadmt_test

import (
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
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
	nb := schemadmt.Type.Schema.Representation().NewBuilder()
	f, err := os.Open("../../.ipld/specs/schemas/schema-schema.ipldsch.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := dagjson.Decode(nb, f); err != nil {
		t.Fatal(err)
	}
	node := nb.Build()
	sch := bindnode.Unwrap(node).(*schemadmt.Schema)
	_ = sch

	// TODO: re-enable testing Compile once it's finished
	var ts schema.TypeSystem
	ts.Init()
	if err := schemadmt.Compile(&ts, sch); err != nil {
		t.Fatal(err)
	}

	typeStruct := ts.TypeByName("TypeStruct")
	println(typeStruct)
}
