package gengo

import (
	"runtime"
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestStructReprStringjoin(t *testing.T) {
	if runtime.GOOS != "darwin" { // TODO: enable parallelism on macos
		t.Parallel()
	}

	engine := &genAndCompileEngine{prefix: "struct-str-join"}
	tests.SchemaTestStructReprStringjoin(t, engine)
}
