package bindnode_test

import (
	"strings"
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
)

// For now, we simply run all schema tests with Prototype.
// In the future, forSchemaTest might return multiple engines.

func forSchemaTest(name string) []tests.EngineSubtest {
	return []tests.EngineSubtest{{
		Engine: &bindEngine{},
	}}
}

func TestSchema(t *testing.T) {
	t.Parallel()

	tests.SchemaTestAll(t, forSchemaTest)
}

var _ tests.Engine = (*bindEngine)(nil)

type bindEngine struct {
	ts schema.TypeSystem
}

func (e *bindEngine) Init(t *testing.T, ts schema.TypeSystem) {
	e.ts = ts
}

func (e *bindEngine) PrototypeByName(name string) ipld.NodePrototype {
	wantRepr := strings.HasSuffix(name, ".Repr")
	if wantRepr {
		name = strings.TrimSuffix(name, ".Repr")
	}
	schemaType := e.ts.TypeByName(name)
	if wantRepr {
		return bindnode.Prototype(nil, schemaType).Representation()
	}
	return bindnode.Prototype(nil, schemaType)
}
