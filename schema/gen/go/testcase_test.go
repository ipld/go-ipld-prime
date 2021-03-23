package gengo

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/shared"
	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/traversal"
)

// This file introduces a testcase struct and a bunch of functions around it.
//  This structure can be used to specify many test scenarios easily, using json as a shorthand for the fixtures.
//  Not everything can be tested this way (in particular, there's some fun details around maps with complex keys, and structs with absent fields), but it covers a lot.

/*
	testcase contains data for directing a sizable number of tests against a NodePrototype
	(or more specifically, a pair of them -- one for the type-level node, one for the representation),
	all of which are applied by calling the testcase.Test method:

		- Creation of values using the type-level builder is tested.
			- This is done using a json input as a convenient shorthand.
			- n.b. this is optional, because it won't work for maps with complex keys.
			- In things that behave as maps: this tests the AssembleEntry path (rather than AssembleKey+AssembleValue; this is the case because this is implemented using unmarshal codepaths).
			- If this is expected to fail, an expected error may be specified (which will also make all other tests after creation inapplicable to this testcase).
		- Creation of values using the repr-level builder is tested.
			- This is (again) done using a json input as a convenient shorthand.
			- At least *one* of this or the json for type-level must be present.  If neither: the testcase spec is broken.
			- As for the type-level test: in things that behave as maps, this tests the AssembleEntry path.
			- If this is expected to fail, an expected error may be specified (which will also make all other tests after creation inapplicable to this testcase).
		- If both forms of creation were exercised: check that the result nodes are deep-equal.
		- A list of "point" observations may be provided, which can probe positions in the data tree for expected values (or just type kind, etc).
			- This tests that direct lookups work.  (It doesn't test iterators; that'll come in another step, later.)
			- Pathing (a la traversal.Get) is used for this this, so it's ready to inspect deep structures.
			- The field for expected value is just `interface{}`; it handles nodes, some primitives, and will also allow asserting an error.
		- The node is *copied*, and deep-equal checked again.
			- The purpose of this is to exercise the AssembleKey+AssembleValue path (as opposed to AssembleEntry (which is already exercised by our creation tests, since they use unmarshal codepaths)).
		- Access of type-level data via iterators is tested in one of two ways:
			- A list of expected key+values expected of the iterator can be provided explicitly;
			- If an explicit list isn't provided, but type-level json is provided, the type-level data will be marshalled and compared to the json fixture.
			- Most things can use the json path -- those that can't (e.g. maps with complex keys; structs with absent values -- neither is marshallable) use the explicit key+value system instead.
		- Access of the representation-level data via interators is tested via marshalling, and asserting it against the json fixture data (if present).
			- There's no explicit key+value list alternative here -- it's not needed; there is no data that is unmarshallable, by design!

	This system should cover a lot of things, but doesn't cover everything.

		- Good coverage for "reset" pathways is reached somewhat indirectly...
			- Tests for recursive types containing nontrivial reset methods exercise both the child type's assembler reset method, and that the parent calls it correctly.
		- Maps with complex keys are tricky to handle, as already noted above.
			- But you should be able to do it, with some care.
		- This whole system depends on json parsers and serializers already working.
			- This is arguably an uncomfortably large and complex dependency for a test system.  However, the json systems are tested by using basicnode; there's no cycle here.
		- "Unhappy paths" in creation are a bit tricky to test.
			- It can be done, but for map-like things, only for the AssembleEntry path.
			- PRs welcome if someone's got a clever idea for a good way to exercise AssembleKey+AssembleValue.  (A variant of unmarshaller implementation?  Would do it; just verbose.)
		- No support yet for checking properties like Length.
			- Future: we could add another type-hinted special case to the testcasePoint.expect for this, i suppose.
*/
type testcase struct {
	name                string          // name for the testcase.
	typeJson            string          // json that will be fed to unmarshal together with a type-level assembler.  marshal output will also be checked for equality.  may be absent.
	reprJson            string          // json that will be fed to unmarshal together with a representational assembler.  marshal output will also be checked for equality.
	expectUnmarshalFail error           // if present, this error will be expected from the unmarshal process (and implicitly, marshal tests will not be applicable for this testcase).
	typePoints          []testcasePoint // inspections that will be made by traversing the type-level nodes.
	reprPoints          []testcasePoint // inspections that will be made by traversing the representation nodes.
	typeItr             []entry         // if set, the type will be iterated in this way.  The remarshalling and checking against typeJson will not be tested.  This is used to probe for correct iteration over Absent values in structs (which needs special handling, because they are unserializable).
	// there's really no need for an 'expectFail' that applies to marshal, because it shouldn't be possible to create data that's unmarshallable!  (excepting data which is not marshallable by some *codec* due to incompleteness of that codec.  But that's not what we're testing, here.)
	// there's no need for a reprItr because the marshalling to reprJson always covers that; unlike with the type level, neither absents nor complex keys can throw a wrench in serialization, so it's always available to us to exercise the iteration code.
}
type testcasePoint struct {
	path   string
	expect interface{} // if primitive: we'll AsFoo and assert equal on that; if an error, we'll expect an error and compare error types; if a kind, we'll check that the thing reached simply has that kind.
}
type entry struct {
	key   interface{} // (mostly string.  not yet defined how this will handle maps with complex keys.)
	value interface{} // same rules as testcasePoint.expect
}

func (tcase testcase) Test(t *testing.T, np, npr ipld.NodePrototype) {
	t.Run(tcase.name, func(t *testing.T) {
		// We'll produce either one or two nodes, depending on the fixture; if two, we'll be expecting them to be equal.
		var n, n2 ipld.Node

		// Attempt to produce a node by using unmarshal on type-level fixture data and the type-level NodePrototype.
		//  This exercises creating a value using the AssembleEntry path (but note, not AssembleKey+AssembleValue path).
		//  This test section is optional because we can't use it for some types (namely, maps with complex keys -- which simply need custom tests).
		if tcase.typeJson != "" {
			t.Run("typed-create", func(t *testing.T) {
				n = testUnmarshal(t, np, tcase.typeJson, tcase.expectUnmarshalFail)
			})
		}

		// Attempt to produce a node by using unmarshal on repr-level fixture data and the repr-level NodePrototype.
		//  This exercises creating a value using the AssembleEntry path (but note, not AssembleKey+AssembleValue path).
		//  This test section is optional simply because it's nice to be able to omit it when writing a new system and not wanting to test representation yet.
		if tcase.reprJson != "" {
			t.Run("repr-create", func(t *testing.T) {
				n3 := testUnmarshal(t, npr, tcase.reprJson, tcase.expectUnmarshalFail)
				if n == nil {
					n = n3
				} else {
					n2 = n3
				}
			})
		}

		// If unmarshalling was expected to fail, the rest of the tests are inapplicable.
		if tcase.expectUnmarshalFail != nil {
			return
		}

		// Check the nodes are equal, if there's two of them.  (Or holler, if none.)
		if n == nil {
			t.Fatalf("invalid fixture: need one of either typeJson or reprJson provided")
		}
		if n2 != nil {
			t.Run("type-create and repr-create match", func(t *testing.T) {
				Wish(t, n, ShouldEqual, n2)
			})
		}

		// Perform all the point inspections on the type-level node.
		if tcase.typePoints != nil {
			t.Run("type-level inspection", func(t *testing.T) {
				for _, point := range tcase.typePoints {
					wishPoint(t, n, point)
				}
			})
		}

		// Perform all the point inspections on the repr-level node.
		if tcase.reprPoints != nil {
			t.Run("repr-level inspection", func(t *testing.T) {
				for _, point := range tcase.reprPoints {
					wishPoint(t, n.(schema.TypedNode).Representation(), point)
				}
			})
		}

		// Copy the node.  This exercises the AssembleKey+AssembleValue path for maps, as opposed to the AssembleEntry path (which was exercised by the creation via unmarshal).
		//  This isn't especially informative for anything other than maps, but we do it universally anyway (it's not a significant time cost).
		// TODO

		// Copy the node, now at repr level.  Again, this is for exercising AssembleKey+AssembleValue paths.
		// TODO

		// Serialize the type-level node, and check that we get the original json again.
		//  This exercises iterators on the type-level node.
		//  OR, if typeItr is present, do that instead (this is necessary when handling maps with complex keys or handling structs with absent values, since both of those are unserializable).
		if tcase.typeItr != nil {
			// This can unconditionally assume we're going to handle maps,
			//  because the only kind of thing that needs this style of testing are some instances of maps and some instances of structs.
			itr := n.MapIterator()
			for _, entry := range tcase.typeItr {
				Wish(t, itr.Done(), ShouldEqual, false)
				k, v, err := itr.Next()
				Wish(t, k, closeEnough, entry.key)
				Wish(t, v, closeEnough, entry.value)
				Wish(t, err, ShouldEqual, nil)
			}
			Wish(t, itr.Done(), ShouldEqual, true)
			k, v, err := itr.Next()
			Wish(t, k, ShouldEqual, nil)
			Wish(t, v, ShouldEqual, nil)
			Wish(t, err, ShouldEqual, ipld.ErrIteratorOverread{})
		} else if tcase.typeJson != "" {
			t.Run("type-marshal", func(t *testing.T) {
				testMarshal(t, n, tcase.typeJson)
			})
		}

		// Serialize the repr-level node, and check that we get the original json again.
		//  This exercises iterators on the repr-level node.
		if tcase.reprJson != "" {
			t.Run("repr-marshal", func(t *testing.T) {
				testMarshal(t, n.(schema.TypedNode).Representation(), tcase.reprJson)
			})
		}
	})

}

func testUnmarshal(t *testing.T, np ipld.NodePrototype, data string, expectFail error) ipld.Node {
	t.Helper()
	nb := np.NewBuilder()
	err := dagjson.Unmarshal(nb, json.NewDecoder(strings.NewReader(data)))
	switch {
	case expectFail == nil && err != nil:
		t.Fatalf("fixture parse failed: %s", err)
	case expectFail == nil && err == nil:
		// carry on
	case expectFail != nil && err != nil:
		Wish(t, err, ShouldBeSameTypeAs, expectFail)
	case expectFail != nil && err == nil:
		t.Errorf("expected creation to fail with a %T error, but got no error", expectFail)
	}
	return nb.Build()
}

func testMarshal(t *testing.T, n ipld.Node, data string) {
	t.Helper()
	// We'll marshal with "pretty" linebreaks and indents (and re-format the fixture to the same) for better diffing.
	prettyprint := json.EncodeOptions{Line: []byte{'\n'}, Indent: []byte{'\t'}}
	var buf bytes.Buffer
	err := dagjson.Marshal(n, json.NewEncoder(&buf, prettyprint), true)
	if err != nil {
		t.Errorf("marshal failed: %s", err)
	}
	Wish(t, buf.String(), ShouldEqual, reformat(data, prettyprint))
}

func wishPoint(t *testing.T, n ipld.Node, point testcasePoint) {
	t.Helper()
	reached, err := traversal.Get(n, ipld.ParsePath(point.path))
	switch point.expect.(type) {
	case error:
		Wish(t, err, ShouldBeSameTypeAs, point.expect)
		Wish(t, err, ShouldEqual, point.expect)
	default:
		Wish(t, err, ShouldEqual, nil)
		if reached == nil {
			return
		}
		Wish(t, reached, closeEnough, point.expect)
	}
}

// closeEnough conforms to wish.Checker (so we can use it in Wish invocations),
// and lets Nodes be compared to primitives in convenient ways.
//
// If the expected value is a primitive string, it'll AsStrong on the Node; etc.
//
// Using an ipld.Kind value is also possible, which will just check the kind and not the value contents.
//
// If an ipld.Node is the expected value, a full deep ShouldEqual is used as normal.
func closeEnough(actual, expected interface{}) (string, bool) {
	if expected == nil {
		return ShouldEqual(actual, nil)
	}
	a, ok := actual.(ipld.Node)
	if !ok {
		return "this checker only supports checking ipld.Node values", false
	}
	switch expected.(type) {
	case ipld.Kind:
		return ShouldEqual(a.Kind(), expected)
	case string:
		if a.Kind() != ipld.Kind_String {
			return fmt.Sprintf("expected something with kind string, got kind %s", a.Kind()), false
		}
		x, _ := a.AsString()
		return ShouldEqual(x, expected)
	case int:
		if a.Kind() != ipld.Kind_Int {
			return fmt.Sprintf("expected something with kind int, got kind %s", a.Kind()), false
		}
		x, _ := a.AsInt()
		return ShouldEqual(x, expected)
	case ipld.Node:
		return ShouldEqual(actual, expected)
	default:
		return fmt.Sprintf("this checker doesn't support an expected value of type %T", expected), false
	}
}

func reformat(x string, opts json.EncodeOptions) string {
	var buf bytes.Buffer
	if err := (shared.TokenPump{
		json.NewDecoder(strings.NewReader(x)),
		json.NewEncoder(&buf, opts),
	}).Run(); err != nil {
		panic(err)
	}
	return buf.String()
}
