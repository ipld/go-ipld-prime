package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestLinks(t *testing.T) {
	t.Parallel()

	engine := &genAndCompileEngine{prefix: "links"}
	tests.SchemaTestLinks(t, engine)
}
