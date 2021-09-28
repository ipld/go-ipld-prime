package schemadmt_test

import (
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

func TestDecodeDAGJSON(t *testing.T) {
	nb := schemadmt.Type.Schema.Representation().NewBuilder()
	f, err := os.Open("../../.ipld/specs/schemas/schema-schema.ipldsch.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := dagjson.Decode(nb, f); err != nil {
		t.Fatal(err)
	}
	node := nb.Build()
	sch := bindnode.Unwrap(node).(*schemadmt.Schema)

	var ts schema.TypeSystem
	ts.Init()
	if err := schemadmt.Compile(&ts, sch); err != nil {
		t.Fatal(err)
	}

	typeStruct := ts.TypeByName("TypeDefnStruct")
	if typeStruct == nil {
		t.Fatal("TypeStruct not found")
	}
}

