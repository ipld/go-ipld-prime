package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestScalars(t *testing.T) {
	t.Parallel()

	engine := &genAndCompileEngine{prefix: "scalars"}
	tests.SchemaTestScalars(t, engine)
}
