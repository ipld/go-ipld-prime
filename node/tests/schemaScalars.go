package tests

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

func assignValue(am datamodel.NodeAssembler, value interface{}) error {
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
		kind  datamodel.Kind
		value interface{}
	}{
		{"Bool", datamodel.Kind_Bool, true},
		{"Int", datamodel.Kind_Int, int64(23)},
		{"Float", datamodel.Kind_Float, 12.25},
		{"String", datamodel.Kind_String, "foo"},
		{"Bytes", datamodel.Kind_Bytes, []byte("bar")},
	}

	// We test each of the five scalar prototypes in subtests.
	for _, testProto := range tests {

		// For both the regular node and its repr version,
		// getting the right value for the kind should work.
		for _, useRepr := range []bool{false, true} {

			protoName := testProto.name
			if useRepr {
				protoName += ".Repr"
			}
			np := engine.PrototypeByName(protoName)

			// For each prototype, we try assigning all scalar values.
			for _, testAssign := range tests {

				// We try both AssignKind and AssignNode.
				for _, useAssignNode := range []bool{false, true} {
					testName := fmt.Sprintf("%s-Assign%s", protoName, testAssign.name)
					if useAssignNode {
						testName = fmt.Sprintf("%s-AssignNode-%s", protoName, testAssign.name)
					}
					t.Run(testName, func(t *testing.T) {
						nb := np.NewBuilder()

						// Assigning null, a list, or a map, should always fail.
						err := nb.AssignNull()
						qt.Check(t, err, qt.Not(qt.IsNil))
						_, err = nb.BeginMap(-1)
						qt.Check(t, err, qt.Not(qt.IsNil))
						_, err = nb.BeginList(-1)
						qt.Check(t, err, qt.Not(qt.IsNil))

						// Assigning the right value for the kind should succeed.
						if useAssignNode {
							np2 := engine.PrototypeByName(testAssign.name)
							nb2 := np2.NewBuilder()
							qt.Check(t, assignValue(nb2, testAssign.value), qt.IsNil)
							n2 := nb2.Build()

							err = nb.AssignNode(n2)
						} else {
							err = assignValue(nb, testAssign.value)
						}
						if testAssign.kind == testProto.kind {
							qt.Check(t, err, qt.IsNil)
						} else {
							qt.Check(t, err, qt.Not(qt.IsNil))

							// Assign something anyway, just so we can Build later.
							err := assignValue(nb, testProto.value)
							qt.Check(t, err, qt.IsNil)
						}

						n := nb.Build()

						var gotValue interface{}
						err = nil
						switch testAssign.kind {
						case datamodel.Kind_Bool:
							gotValue, err = n.AsBool()
						case datamodel.Kind_Int:
							gotValue, err = n.AsInt()
						case datamodel.Kind_Float:
							gotValue, err = n.AsFloat()
						case datamodel.Kind_String:
							gotValue, err = n.AsString()
						case datamodel.Kind_Bytes:
							gotValue, err = n.AsBytes()
						default:
							t.Fatal(testAssign.kind)
						}
						if testAssign.kind == testProto.kind {
							qt.Check(t, err, qt.IsNil)
							qt.Check(t, gotValue, qt.DeepEquals, testAssign.value)
						} else {
							qt.Check(t, err, qt.Not(qt.IsNil))
						}

						// Using Node methods which should never
						// work on scalar kinds.

						_, err = n.LookupByString("foo")
						qt.Check(t, err, qt.Not(qt.IsNil))
						_, err = n.LookupByIndex(3)
						qt.Check(t, err, qt.Not(qt.IsNil))
						qt.Check(t, n.MapIterator(), qt.IsNil)
						qt.Check(t, n.ListIterator(), qt.IsNil)
						qt.Check(t, n.Length(), qt.Equals, int64(-1))
						qt.Check(t, n.IsAbsent(), qt.IsFalse)
						qt.Check(t, n.IsNull(), qt.IsFalse)
					})
				}
			}
		}
	}
}
