package schemadsl_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/schema"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
	schemadsl "github.com/ipld/go-ipld-prime/schema/dsl"
)

func TestParse(t *testing.T) {
	sch, err := schemadsl.ParseFile("../../.ipld/specs/schemas/schema-schema.ipldsch")
	if err != nil {
		t.Fatal(err)
	}

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
