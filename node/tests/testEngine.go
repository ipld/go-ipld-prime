package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

// Engine describes the interface that can be supplied to run tests on schemas.
//
// The PrototypeByName function can get its job done using only interface types
// that we already know from outside any generated code, so you can write tests
// that have no _compile time_ dependency on the generated code.  This makes it
// easier for IDEs and suchlike to help you write and check the test functions.
//
// Ask for prototypes using the type name alone (no package prefix);
// their representation prototypes can be obtained by appending ".Repr".
type Engine interface {
	Init(t *testing.T, ts schema.TypeSystem)
	PrototypeByName(name string) datamodel.NodePrototype
}
