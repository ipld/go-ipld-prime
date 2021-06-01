package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestString(t *testing.T) {
	t.Parallel()

	engine := &genAndCompileEngine{prefix: "string"}
	tests.SchemaTestString(t, engine)
}
