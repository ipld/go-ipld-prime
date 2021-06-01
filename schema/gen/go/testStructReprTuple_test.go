package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestStructReprTuple(t *testing.T) {
	t.Parallel()

	engine := &genAndCompileEngine{prefix: "struct-tuple"}
	tests.SchemaTestStructReprTuple(t, engine)
}
