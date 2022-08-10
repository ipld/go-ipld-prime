package bindnode_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/printer"
)

// test building a struct with two string fields
func TestBasicStructBuilding(t *testing.T) {
	ts, err := ipld.LoadSchemaBytes([]byte(`
type Animal struct {
	Name  String
	Sound String
}
	`))
	if err != nil {
		t.Fatal(err)
	}

	schemaType := ts.TypeByName("Animal")
	proto := bindnode.Prototype(nil, schemaType)

	b := proto.NewBuilder()
	ma, err := b.BeginMap(2)
	if err != nil {
		t.Fatal(err)
	}

	a := ma.AssembleKey()
	must.NotError(a.AssignString("Name"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("cat"))

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Sound"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("meow"))

	ma.Finish()

	mapNode := b.Build()
	actual := printer.Sprint(mapNode)

	expect := `struct<Animal>{
	Name: string<String>{"cat"}
	Sound: string<String>{"meow"}
}`
	qt.Check(t, expect, qt.Equals, actual)

	qt.Check(t, fmt.Sprintf("%T", proto.Representation()), qt.Equals, "*bindnode._prototypeRepr")
	qt.Check(t, fmt.Sprintf("%s", schemaType.TypeKind()), qt.Equals, "struct")
	qt.Check(t, fmt.Sprintf("%s", schemaType.RepresentationBehavior()), qt.Equals, "map")
}

// test building a struct that contains a union, using typed nodes
func TestStructWithUnion(t *testing.T) {
	ts, err := ipld.LoadSchemaBytes([]byte(`
type Animal struct {
	Name   String
	Action Behavior
}

type Behavior union {
	| Movement "movement"
	| Sound    "sound"
} representation keyed

type Movement string

type Sound struct {
	Vocal String
	Alt   String
}
	`))
	if err != nil {
		t.Fatal(err)
	}

	// Sound struct with two keys
	schemaType := ts.TypeByName("Sound")
	proto := bindnode.Prototype(nil, schemaType)

	b := proto.NewBuilder()
	ma, err := b.BeginMap(2)
	if err != nil {
		t.Fatal(err)
	}

	a := ma.AssembleKey()
	must.NotError(a.AssignString("Vocal"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("bark"))

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Alt"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("woof"))

	ma.Finish()
	soundNode := b.Build()

	// Behavior union, contains Sound
	schemaType = ts.TypeByName("Behavior")
	proto = bindnode.Prototype(nil, schemaType)

	b = proto.NewBuilder()
	ma, err = b.BeginMap(1)
	if err != nil {
		t.Fatal(err)
	}

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Sound"))
	a = ma.AssembleValue()
	must.NotError(a.AssignNode(soundNode))

	ma.Finish()
	behaviorNode := b.Build()

	// Animal top-level struct
	schemaType = ts.TypeByName("Animal")
	proto = bindnode.Prototype(nil, schemaType)

	b = proto.NewBuilder()
	ma, err = b.BeginMap(2)
	if err != nil {
		t.Fatal(err)
	}

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Name"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("dog"))

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Action"))
	a = ma.AssembleValue()
	must.NotError(a.AssignNode(behaviorNode))

	ma.Finish()
	animalNode := b.Build()
	actual := printer.Sprint(animalNode)

	expect := `struct<Animal>{
	Name: string<String>{"dog"}
	Action: union<Behavior>{struct<Sound>{
		Vocal: string<String>{"bark"}
		Alt: string<String>{"woof"}
	}}
}`
	qt.Check(t, expect, qt.Equals, actual)

	// different Behavior for a different animal
	schemaType = ts.TypeByName("Behavior")
	proto = bindnode.Prototype(nil, schemaType)

	b = proto.NewBuilder()
	ma, err = b.BeginMap(1)
	if err != nil {
		t.Fatal(err)
	}

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Movement"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("swim"))

	ma.Finish()
	behaviorNode = b.Build()

	// another animal, using Movement instead of Sound
	schemaType = ts.TypeByName("Animal")
	proto = bindnode.Prototype(nil, schemaType)

	b = proto.NewBuilder()
	ma, err = b.BeginMap(2)
	if err != nil {
		t.Fatal(err)
	}

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Name"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("eel"))

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Action"))
	a = ma.AssembleValue()
	must.NotError(a.AssignNode(behaviorNode))

	ma.Finish()
	animalNode = b.Build()
	actual = printer.Sprint(animalNode)

	expect = `struct<Animal>{
	Name: string<String>{"eel"}
	Action: union<Behavior>{string<Movement>{"swim"}}
}`
	qt.Check(t, expect, qt.Equals, actual)

}

// test building a struct that contains a union, where the union is
// constructed via representation
func TestStructWithUnionByRepresentation(t *testing.T) {
	ts, err := ipld.LoadSchemaBytes([]byte(`
type Animal struct {
	Name   String
	Action Behavior
}

type Behavior union {
	| Movement "movement:"
	| Sound    "sound:"
} representation stringprefix

type Movement string

type Sound struct {
	Vocal String
	Alt   String
}
	`))
	if err != nil {
		t.Fatal(err)
	}

	// Behavior union can't be built from representation using this prototype
	// attempting to do so will cause an error
	schemaType := ts.TypeByName("Behavior")
	typedProto := bindnode.Prototype(nil, schemaType)

	b := typedProto.NewBuilder()
	err = b.AssignString("movement:swim")
	qt.Check(t, err, qt.ErrorMatches, `func called on wrong kind: "AssignString" called on a Behavior node \(kind: map\), but only makes sense on string`)

	// Behavior union can be built using representation if requested
	schemaType = ts.TypeByName("Behavior")
	proto := bindnode.Prototype(nil, schemaType).Representation()

	b = proto.NewBuilder()
	must.NotError(b.AssignString("movement:swim"))
	behaviorNode := b.Build()

	actual := printer.Sprint(behaviorNode)

	expect := `union<Behavior>{string<Movement>{"swim"}}`
	qt.Check(t, expect, qt.Equals, actual)

	// Animal top-level struct
	schemaType = ts.TypeByName("Animal")
	proto = bindnode.Prototype(nil, schemaType)

	b = proto.NewBuilder()
	ma, err := b.BeginMap(2)
	if err != nil {
		t.Fatal(err)
	}

	a := ma.AssembleKey()
	must.NotError(a.AssignString("Name"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("eel"))

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Action"))
	a = ma.AssembleValue()
	must.NotError(a.AssignNode(behaviorNode))

	ma.Finish()
	animalNode := b.Build()
	actual = printer.Sprint(animalNode)

	expect = `struct<Animal>{
	Name: string<String>{"eel"}
	Action: union<Behavior>{string<Movement>{"swim"}}
}`
	qt.Check(t, expect, qt.Equals, actual)
}

// test building a struct that contains a union, where the union is
// constructed as a typed node, but the struct is constructed using
// the representation
func TestStructByRepresentationWithUnion(t *testing.T) {
	ts, err := ipld.LoadSchemaBytes([]byte(`
type Animal struct {
	Name   String
	Action Behavior
}

type Behavior union {
	| Movement "movement:"
	| Sound    "sound:"
} representation stringprefix

type Movement string

type Sound struct {
	Vocal String
	Alt   String
}
	`))
	if err != nil {
		t.Fatal(err)
	}

	// Behavior union built as normal
	schemaType := ts.TypeByName("Behavior")
	typedProto := bindnode.Prototype(nil, schemaType)

	b := typedProto.NewBuilder()
	ma, err := b.BeginMap(1)
	if err != nil {
		t.Fatal(err)
	}

	a := ma.AssembleKey()
	must.NotError(a.AssignString("Movement"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("swim"))
	must.NotError(ma.Finish())
	behaviorNode := b.Build()

	actual := printer.Sprint(behaviorNode)

	expect := `union<Behavior>{string<Movement>{"swim"}}`
	qt.Check(t, expect, qt.Equals, actual)

	// Animal top-level, from representation
	// NOTE: the Behavior node was built as a typed node, but
	// this Animmal assembler uses representation
	schemaType = ts.TypeByName("Animal")
	proto := bindnode.Prototype(nil, schemaType).Representation()

	b = proto.NewBuilder()
	ma, err = b.BeginMap(2)
	if err != nil {
		t.Fatal(err)
	}

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Name"))
	a = ma.AssembleValue()
	must.NotError(a.AssignString("eel"))

	a = ma.AssembleKey()
	must.NotError(a.AssignString("Action"))
	a = ma.AssembleValue()
	must.NotError(a.AssignNode(behaviorNode))

	ma.Finish()
	animalNode := b.Build()
	actual = printer.Sprint(animalNode)

	expect = `struct<Animal>{
	Name: string<String>{"eel"}
	Action: union<Behavior>{string<Movement>{"swim"}}
}`
	qt.Check(t, expect, qt.Equals, actual)
}
