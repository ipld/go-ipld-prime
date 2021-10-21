package bindnode_test

import (
	"strings"
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
)

// For now, we simply run all schema tests with Prototype.
// In the future, forSchemaTest might return multiple engines.

func forSchemaTest(name string) []tests.EngineSubtest {
	if name == "Links" {
		// TODO(mvdan): support typed links; see https://github.com/ipld/go-ipld-prime/issues/272
		return nil
	}
	if name == "UnionKeyedComplexChildren" {
		return nil // Specifically, 'InhabitantB/repr-create_with_AK+AV' borks, because it needs representation-level AssignNode to support more.
	}
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

func (e *bindEngine) PrototypeByName(name string) datamodel.NodePrototype {
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
