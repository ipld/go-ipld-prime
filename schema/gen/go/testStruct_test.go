package gengo

import (
	"runtime"
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestRequiredFields(t *testing.T) {
	if runtime.GOOS != "darwin" { // TODO: enable parallelism on macos
		t.Parallel()
	}

	engine := &genAndCompileEngine{prefix: "struct-required-fields"}
	tests.SchemaTestRequiredFields(t, engine)
}
