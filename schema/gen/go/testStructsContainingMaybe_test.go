package gengo

import (
	"runtime"
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestStructsContainingMaybe(t *testing.T) {
	if runtime.GOOS != "darwin" { // TODO: enable parallelism on macos
		t.Parallel()
	}

	for _, engine := range []*genAndCompileEngine{
		{
			subtestName: "maybe-using-embed",
			prefix:      "stroct",
			adjCfg: AdjunctCfg{
				maybeUsesPtr: map[schema.TypeName]bool{"String": false},
			},
		},
		{
			subtestName: "maybe-using-ptr",
			prefix:      "stroct2",
			adjCfg: AdjunctCfg{
				maybeUsesPtr: map[schema.TypeName]bool{"String": false},
			},
		},
	} {
		t.Run(engine.subtestName, func(t *testing.T) {
			tests.SchemaTestStructsContainingMaybe(t, engine)
		})
	}
}
