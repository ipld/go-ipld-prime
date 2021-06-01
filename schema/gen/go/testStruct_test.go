package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestRequiredFields(t *testing.T) {
	t.Parallel()

	engine := &genAndCompileEngine{prefix: "struct-required-fields"}
	tests.SchemaTestRequiredFields(t, engine)
}
