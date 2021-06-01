package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestStructReprStringjoin(t *testing.T) {
	t.Parallel()

	engine := &genAndCompileEngine{prefix: "struct-str-join"}
	tests.SchemaTestStructReprStringjoin(t, engine)
}
