package basicnode_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/printer"
)

func TestMap(t *testing.T) {
	tests.SpecTestMapStrInt(t, basicnode.Prototype.Map)
	tests.SpecTestMapStrMapStrInt(t, basicnode.Prototype.Map)
	tests.SpecTestMapStrListStr(t, basicnode.Prototype.Map)
}

func BenchmarkMapStrInt_3n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleStandard(b, basicnode.Prototype.Map)
}
func BenchmarkMapStrInt_3n_AssembleEntry(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleEntry(b, basicnode.Prototype.Map)
}
func BenchmarkMapStrInt_3n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_Iteration(b, basicnode.Prototype.Map)
}

func BenchmarkMapStrInt_25n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleStandard(b, basicnode.Prototype.Map)
}
func BenchmarkMapStrInt_25n_AssembleEntry(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleEntry(b, basicnode.Prototype.Map)
}
func BenchmarkMapStrInt_25n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_Iteration(b, basicnode.Prototype.Map)
}

func BenchmarkSpec_Marshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt(b, basicnode.Prototype.Map)
}
func BenchmarkSpec_Marshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_MapNStrMap3StrInt(b, basicnode.Prototype.Map)
}

func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_Map3StrInt(b, basicnode.Prototype.Map)
}
func BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b, basicnode.Prototype.Map)
}

// Test that the map builder cannot be assigned arbitrary values, and trying to
// will result in a sensible error
func TestMapAssignError(t *testing.T) {
	b := basicnode.Prototype.Map.NewBuilder()
	err := b.AssignBool(true)
	errExpect := `func called on wrong kind: "AssignBool" called on a map node \(kind: map\), but only makes sense on bool`
	qt.Check(t, err, qt.ErrorMatches, errExpect)

	err = b.AssignInt(3)
	errExpect = `func called on wrong kind: "AssignInt" called on a map node \(kind: map\), but only makes sense on int`
	qt.Check(t, err, qt.ErrorMatches, errExpect)

	err = b.AssignFloat(5.7)
	errExpect = `func called on wrong kind: "AssignFloat" called on a map node \(kind: map\), but only makes sense on float`
	qt.Check(t, err, qt.ErrorMatches, errExpect)

	err = b.AssignString("hi")
	errExpect = `func called on wrong kind: "AssignString" called on a map node \(kind: map\), but only makes sense on string`
	qt.Check(t, err, qt.ErrorMatches, errExpect)

	err = b.AssignNode(basicnode.NewInt(3))
	errExpect = `func called on wrong kind: "AssignNode" called on a map node \(kind: int\), but only makes sense on map`
	qt.Check(t, err, qt.ErrorMatches, errExpect)

	// TODO(dustmop): BeginList, AssignNull, AssignBytes, AssignLink
}

// Test that the map builder can create map nodes, and AssignNode will copy
// such a node, and lookup methods on that node will work correctly
func TestMapBuilder(t *testing.T) {
	b := basicnode.Prototype.Map.NewBuilder()

	// construct a map of three keys, using the MapBuilder
	ma, err := b.BeginMap(3)
	if err != nil {
		t.Fatal(err)
	}
	a := ma.AssembleKey()
	a.AssignString("cat")
	a = ma.AssembleValue()
	a.AssignString("meow")

	a, err = ma.AssembleEntry("dog")
	if err != nil {
		t.Fatal(err)
	}
	a.AssignString("bark")

	a = ma.AssembleKey()
	a.AssignString("eel")
	a = ma.AssembleValue()
	a.AssignString("zap")

	err = ma.Finish()
	if err != nil {
		t.Fatal(err)
	}

	// test the builder's prototypes and its key and value prototypes, while we're here
	np := b.Prototype()
	qt.Check(t, fmt.Sprintf("%T", np), qt.Equals, "basicnode.Prototype__Map")
	np = ma.KeyPrototype()
	qt.Check(t, fmt.Sprintf("%T", np), qt.Equals, "basicnode.Prototype__String")
	np = ma.ValuePrototype("")
	qt.Check(t, fmt.Sprintf("%T", np), qt.Equals, "basicnode.Prototype__Any")

	// compare the printed map
	mapNode := b.Build()
	actual := printer.Sprint(mapNode)

	expect := `map{
	string{"cat"}: string{"meow"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`
	qt.Check(t, expect, qt.Equals, actual)

	// copy the map using AssignNode
	c := basicnode.Prototype.Map.NewBuilder()
	err = c.AssignNode(mapNode)
	if err != nil {
		t.Fatal(err)
	}
	anotherNode := c.Build()

	actual = printer.Sprint(anotherNode)
	qt.Assert(t, expect, qt.Equals, actual)

	// access values of map, using string
	r, err := anotherNode.LookupByString("cat")
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "meow", qt.Equals, must.String(r))

	// access values of map, using node
	r, err = anotherNode.LookupByNode(basicnode.NewString("dog"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "bark", qt.Equals, must.String(r))

	// access values of map, using PathSegment
	r, err = anotherNode.LookupBySegment(datamodel.ParsePathSegment("eel"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "zap", qt.Equals, must.String(r))

	// validate the node's prototype
	np = anotherNode.Prototype()
	qt.Check(t, fmt.Sprintf("%T", np), qt.Equals, "basicnode.Prototype__Map")
}

// test that AssignNode will fail if called twice, it expects an empty
// node to assign to
func TestMapCantAssignNodeTwice(t *testing.T) {
	b := basicnode.Prototype.Map.NewBuilder()

	// construct a map of three keys, using the MapBuilder
	ma, err := b.BeginMap(3)
	if err != nil {
		t.Fatal(err)
	}
	a := ma.AssembleKey()
	a.AssignString("cat")
	a = ma.AssembleValue()
	a.AssignString("meow")

	a, err = ma.AssembleEntry("dog")
	if err != nil {
		t.Fatal(err)
	}
	a.AssignString("bark")

	a = ma.AssembleKey()
	a.AssignString("eel")
	a = ma.AssembleValue()
	a.AssignString("zap")

	err = ma.Finish()
	if err != nil {
		t.Fatal(err)
	}
	mapNode := b.Build()

	// copy the map using AssignNode, works the first time
	c := basicnode.Prototype.Map.NewBuilder()
	err = c.AssignNode(mapNode)
	if err != nil {
		t.Fatal(err)
	}
	qt.Assert(t,
		func() {
			_ = c.AssignNode(mapNode)
		},
		qt.PanicMatches,
		// TODO(dustmop): Error message here should be better
		`misuse`)
}

func TestMapLookupError(t *testing.T) {
	b := basicnode.Prototype.Map.NewBuilder()

	// construct a map of three keys, using the MapBuilder
	ma, err := b.BeginMap(3)
	if err != nil {
		t.Fatal(err)
	}
	a := ma.AssembleKey()
	a.AssignString("cat")
	a = ma.AssembleValue()
	a.AssignString("meow")

	a, err = ma.AssembleEntry("dog")
	if err != nil {
		t.Fatal(err)
	}
	a.AssignString("bark")

	a = ma.AssembleKey()
	a.AssignString("eel")
	a = ma.AssembleValue()
	a.AssignString("zap")

	err = ma.Finish()
	if err != nil {
		t.Fatal(err)
	}

	mapNode := b.Build()

	_, err = mapNode.LookupByString("frog")
	qt.Check(t, err, qt.ErrorMatches, `key not found: "frog"`)

	_, err = mapNode.LookupByNode(basicnode.NewInt(3))
	// TODO(dustmop): This error message is not great. It's about how the
	// int node could not be converted when the real problem is that this
	// method should not accept ints as parameters
	qt.Check(t, err, qt.ErrorMatches, `func called on wrong kind: "AsString" called on a int node \(kind: int\), but only makes sense on string`)

	_, err = mapNode.LookupByIndex(0)
	qt.Check(t, err, qt.ErrorMatches, `func called on wrong kind: "LookupByIndex" called on a map node \(kind: map\), but only makes sense on list`)
}

func TestMapNewBuilderUsageError(t *testing.T) {
	// construct an empty map
	b := basicnode.Prototype.Map.NewBuilder()
	ma, err := b.BeginMap(0)
	if err != nil {
		t.Fatal(err)
	}
	err = ma.Finish()
	if err != nil {
		t.Fatal(err)
	}
	mapNode := b.Build()
	actual := printer.Sprint(mapNode)

	expect := `map{}`
	qt.Check(t, expect, qt.Equals, actual)

	// assembling a key without a value will cause Finish to panic
	b.Reset()
	ma, err = b.BeginMap(0)
	if err != nil {
		t.Fatal(err)
	}
	a := ma.AssembleKey()
	a.AssignString("cat")
	qt.Assert(t,
		func() {
			_ = ma.Finish()
		},
		qt.PanicMatches,
		// TODO(dustmop): Error message here should be better
		`misuse`)
}

func TestMapDupKeyError(t *testing.T) {
	b := basicnode.Prototype.Map.NewBuilder()

	// construct a map with duplicate keys
	ma, err := b.BeginMap(3)
	if err != nil {
		t.Fatal(err)
	}
	a := ma.AssembleKey()
	a.AssignString("cat")
	a = ma.AssembleValue()
	a.AssignString("meow")
	a = ma.AssembleKey()
	err = a.AssignString("cat")

	qt.Check(t, err, qt.ErrorMatches, `cannot repeat map key "cat"`)
}

func TestMapAmendingBuilderNewNode(t *testing.T) {
	// Create a map amender with an empty base node
	amender := basicnode.Prototype.Map.AmendingBuilder(nil)

	err := amender.Put(basicnode.NewString("cat"), basicnode.NewString("meow"))
	if err != nil {
		t.Fatal(err)
	}
	// Retry adding the entry
	err = amender.Put(basicnode.NewString("cat"), basicnode.NewString("meow"))
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Put(basicnode.NewString("dog"), basicnode.NewString("bark"))
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Put(basicnode.NewString("eel"), basicnode.NewString("zap"))
	if err != nil {
		t.Fatal(err)
	}

	// compare the printed map
	mapNode := amender.Build()
	actual := printer.Sprint(mapNode)

	expect := `map{
	string{"cat"}: string{"meow"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`
	qt.Check(t, expect, qt.Equals, actual)

	// Access values of map, using string
	r, err := mapNode.LookupByString("cat")
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "meow", qt.Equals, must.String(r))

	// Access values of map, using node
	r, err = mapNode.LookupByNode(basicnode.NewString("dog"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "bark", qt.Equals, must.String(r))

	// Access values of map, using PathSegment
	r, err = mapNode.LookupBySegment(datamodel.ParsePathSegment("eel"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "zap", qt.Equals, must.String(r))

	// Validate the node's prototype
	np := mapNode.Prototype()
	qt.Check(t, fmt.Sprintf("%T", np), qt.Equals, "basicnode.Prototype__Map")

	// Amend the map
	err = amender.Put(basicnode.NewString("cat"), basicnode.NewString("purr"))
	if err != nil {
		t.Fatal(err)
	}

	// Access updated value of map, using get
	r, err = amender.Get(basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "purr", qt.Equals, must.String(r))

	expect = `map{
	string{"cat"}: string{"purr"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`

	// The original node should have been updated
	actual = printer.Sprint(mapNode)
	qt.Assert(t, expect, qt.Equals, actual)

	// Remove an entry
	removed, err := amender.Remove(basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Assert(t, removed, qt.IsTrue, qt.Commentf("remove should have returned true"))

	expect = `map{
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`

	// The original node should have been updated
	actual = printer.Sprint(mapNode)
	qt.Assert(t, expect, qt.Equals, actual)

	// Should not find "cat"
	r, err = amender.Get(basicnode.NewString("cat"))
	if _, notFoundErr := err.(datamodel.ErrNotExists); !notFoundErr {
		t.Fatal(err)
	}

	keys, err := amender.Keys()
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"dog"}
	1: string{"eel"}
}`
	actual = printer.Sprint(keys)
	qt.Assert(t, expect, qt.Equals, actual)

	values, err := amender.Values()
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"bark"}
	1: string{"zap"}
}`
	actual = printer.Sprint(values)
	qt.Assert(t, expect, qt.Equals, actual)
}

func TestMapAmendingBuilderExistingNode(t *testing.T) {
	b := basicnode.Prototype.Map.NewBuilder()

	// construct a map of three keys, using the MapBuilder
	ma, err := b.BeginMap(3)
	if err != nil {
		t.Fatal(err)
	}
	a := ma.AssembleKey()
	a.AssignString("cat")
	a = ma.AssembleValue()
	a.AssignString("meow")

	a, err = ma.AssembleEntry("dog")
	if err != nil {
		t.Fatal(err)
	}
	a.AssignString("bark")

	a = ma.AssembleKey()
	a.AssignString("eel")
	a = ma.AssembleValue()
	a.AssignString("zap")

	err = ma.Finish()
	if err != nil {
		t.Fatal(err)
	}

	// Wrap in an amending builder
	amender := basicnode.Prototype.Map.AmendingBuilder(b.Build())

	// Compare the printed map
	mapNode := b.Build()
	actual := printer.Sprint(mapNode)

	expect := `map{
	string{"cat"}: string{"meow"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`
	qt.Check(t, expect, qt.Equals, actual)

	// Access values of map, using string
	r, err := mapNode.LookupByString("cat")
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "meow", qt.Equals, must.String(r))

	// Access values of map, using node
	r, err = mapNode.LookupByNode(basicnode.NewString("dog"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "bark", qt.Equals, must.String(r))

	// Access values of map, using PathSegment
	r, err = mapNode.LookupBySegment(datamodel.ParsePathSegment("eel"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "zap", qt.Equals, must.String(r))

	// Validate the node's prototype
	np := mapNode.Prototype()
	qt.Check(t, fmt.Sprintf("%T", np), qt.Equals, "basicnode.Prototype__Map")

	// Amend the map
	err = amender.Put(basicnode.NewString("cat"), basicnode.NewString("purr"))
	if err != nil {
		t.Fatal(err)
	}

	// Access updated value of map, using get
	r, err = amender.Get(basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "purr", qt.Equals, must.String(r))

	expect = `map{
	string{"cat"}: string{"purr"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`

	// The original node should have been updated
	actual = printer.Sprint(mapNode)
	qt.Assert(t, expect, qt.Equals, actual)

	// Remove an entry
	removed, err := amender.Remove(basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Assert(t, removed, qt.IsTrue, qt.Commentf("remove should have returned true"))

	expect = `map{
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`

	// The original node should have been updated
	actual = printer.Sprint(mapNode)
	qt.Assert(t, expect, qt.Equals, actual)

	// Should not find "cat"
	r, err = amender.Get(basicnode.NewString("cat"))
	if _, notFoundErr := err.(datamodel.ErrNotExists); !notFoundErr {
		t.Fatal(err)
	}

	keys, err := amender.Keys()
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"dog"}
	1: string{"eel"}
}`
	actual = printer.Sprint(keys)
	qt.Assert(t, expect, qt.Equals, actual)

	values, err := amender.Values()
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"bark"}
	1: string{"zap"}
}`
	actual = printer.Sprint(values)
	qt.Assert(t, expect, qt.Equals, actual)
}

func TestMapAmendingBuilderCopiedNode(t *testing.T) {
	b := basicnode.Prototype.Map.NewBuilder()

	// Construct a map of three keys, using the MapBuilder
	ma, err := b.BeginMap(3)
	if err != nil {
		t.Fatal(err)
	}
	a := ma.AssembleKey()
	a.AssignString("cat")
	a = ma.AssembleValue()
	a.AssignString("meow")

	a, err = ma.AssembleEntry("dog")
	if err != nil {
		t.Fatal(err)
	}
	a.AssignString("bark")

	a = ma.AssembleKey()
	a.AssignString("eel")
	a = ma.AssembleValue()
	a.AssignString("zap")

	err = ma.Finish()
	if err != nil {
		t.Fatal(err)
	}

	origMapNode := b.Build()
	// Copy the map using datamodel.Copy. AssignNode will copy pointers to internal values and will not return a fully
	// standalone copy.
	origAmender := basicnode.Prototype.Map.AmendingBuilder(nil)
	err = datamodel.Copy(origMapNode, origAmender)
	if err != nil {
		t.Fatal(err)
	}

	// Wrap in an amending builder
	newAmender := basicnode.Prototype.Map.AmendingBuilder(origAmender.Build())

	// Compare the printed map
	newMapNode := newAmender.Build()
	actual := printer.Sprint(newMapNode)

	expect := `map{
	string{"cat"}: string{"meow"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`
	qt.Check(t, expect, qt.Equals, actual)

	// Amend the copied map
	err = newAmender.Put(basicnode.NewString("cat"), basicnode.NewString("purr"))
	if err != nil {
		t.Fatal(err)
	}

	expect = `map{
	string{"cat"}: string{"purr"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`

	// The new node should have been updated
	actual = printer.Sprint(newMapNode)
	qt.Assert(t, expect, qt.Equals, actual)

	// Remove an entry
	removed, err := newAmender.Remove(basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Assert(t, removed, qt.IsTrue, qt.Commentf("remove should have returned true"))

	expect = `map{
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`

	_, err = newAmender.Get(basicnode.NewString("cat"))
	if _, notFoundErr := err.(datamodel.ErrNotExists); !notFoundErr {
		t.Fatal(err)
	}

	keys, err := newAmender.Keys()
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"dog"}
	1: string{"eel"}
}`
	actual = printer.Sprint(keys)
	qt.Assert(t, expect, qt.Equals, actual)

	values, err := newAmender.Values()
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"bark"}
	1: string{"zap"}
}`
	actual = printer.Sprint(values)
	qt.Assert(t, expect, qt.Equals, actual)

	// The original node should not have been updated
	expect = `map{
	string{"cat"}: string{"meow"}
	string{"dog"}: string{"bark"}
	string{"eel"}: string{"zap"}
}`
	actual = printer.Sprint(origMapNode)
	qt.Assert(t, expect, qt.Equals, actual)
}
