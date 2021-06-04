package tests

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func assignValue(am ipld.NodeAssembler, value interface{}) error {
	switch value := value.(type) {
	case bool:
		return am.AssignBool(value)
	case int64:
		return am.AssignInt(value)
	case float64:
		return am.AssignFloat(value)
	case string:
		return am.AssignString(value)
	case []byte:
		return am.AssignBytes(value)
	default:
		panic(fmt.Sprintf("%T", value))
	}
}

func SchemaTestScalars(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()

	ts.Accumulate(schema.SpawnBool("Bool"))
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnFloat("Float"))
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnBytes("Bytes"))
	engine.Init(t, ts)

	var tests = []struct {
		name  string
		kind  ipld.Kind
		value interface{}
	}{
		{"Bool", ipld.Kind_Bool, true},
		{"Int", ipld.Kind_Int, int64(23)},
		{"Float", ipld.Kind_Float, 12.25},
		{"String", ipld.Kind_String, "foo"},
		{"Bytes", ipld.Kind_Bytes, []byte("bar")},
	}

	// We test each of the five scalar prototypes in subtests.
	for _, testProto := range tests {
		np := engine.PrototypeByName(testProto.name)

		// For each prototype, we try assigning all scalar values.
		for _, testAssign := range tests {

			// We try both AssignKind and AssignNode.
			for _, useAssignNode := range []bool{false, true} {
				testName := fmt.Sprintf("%s-Assign%s", testProto.name, testAssign.name)
				if useAssignNode {
					testName = fmt.Sprintf("%s-AssignNode-%s", testProto.name, testAssign.name)
				}
				t.Run(testName, func(t *testing.T) {
					nb := np.NewBuilder()

					// Assigning null, a list, or a map, should always fail.
					err := nb.AssignNull()
					qt.Assert(t, err, qt.Not(qt.IsNil))
					_, err = nb.BeginMap(-1)
					qt.Assert(t, err, qt.Not(qt.IsNil))
					_, err = nb.BeginList(-1)
					qt.Assert(t, err, qt.Not(qt.IsNil))

					// Assigning the right value for the kind should succeed.
					if useAssignNode {
						np2 := engine.PrototypeByName(testAssign.name)
						nb2 := np2.NewBuilder()
						qt.Assert(t, assignValue(nb2, testAssign.value), qt.IsNil)
						n2 := nb2.Build()

						err = nb.AssignNode(n2)
					} else {
						err = assignValue(nb, testAssign.value)
					}
					if testAssign.kind == testProto.kind {
						qt.Assert(t, err, qt.IsNil)
					} else {
						qt.Assert(t, err, qt.Not(qt.IsNil))

						// Assign something anyway, just so we can Build later.
						err := assignValue(nb, testProto.value)
						qt.Assert(t, err, qt.IsNil)
					}

					n := nb.Build()

					// For both the regular node and its repr version,
					// getting the right value for the kind should work.
					for _, n := range []ipld.Node{
						n,
						n.(schema.TypedNode).Representation(),
					} {
						var gotValue interface{}
						err = nil
						switch testAssign.kind {
						case ipld.Kind_Bool:
							gotValue, err = n.AsBool()
						case ipld.Kind_Int:
							gotValue, err = n.AsInt()
						case ipld.Kind_Float:
							gotValue, err = n.AsFloat()
						case ipld.Kind_String:
							gotValue, err = n.AsString()
						case ipld.Kind_Bytes:
							gotValue, err = n.AsBytes()
						default:
							t.Fatal(testAssign.kind)
						}
						if testAssign.kind == testProto.kind {
							qt.Assert(t, err, qt.IsNil)
							qt.Assert(t, gotValue, qt.DeepEquals, testAssign.value)
						} else {
							qt.Assert(t, err, qt.Not(qt.IsNil))
						}

						// Using Node methods which should never
						// work on scalar kinds.

						_, err = n.LookupByString("foo")
						qt.Assert(t, err, qt.Not(qt.IsNil))
						_, err = n.LookupByIndex(3)
						qt.Assert(t, err, qt.Not(qt.IsNil))
						qt.Assert(t, n.MapIterator(), qt.IsNil)
						qt.Assert(t, n.ListIterator(), qt.IsNil)
						qt.Assert(t, n.Length(), qt.Equals, int64(-1))
						qt.Assert(t, n.IsAbsent(), qt.IsFalse)
						qt.Assert(t, n.IsNull(), qt.IsFalse)
					}
				})
			}
		}
	}
}
