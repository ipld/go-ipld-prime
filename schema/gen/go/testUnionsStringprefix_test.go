package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestUnionStringprefix(t *testing.T) {
	t.Parallel()

	for _, engine := range []*genAndCompileEngine{
		{
			subtestName: "union-using-embed",
			prefix:      "union-stringprefix-using-embed",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"WheeUnion": "embedAll"},
			},
		},
		{
			subtestName: "union-using-interface",
			prefix:      "union-stringprefix-using-interface",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"WheeUnion": "interface"},
			},
		},
	} {
		t.Run(engine.subtestName, func(t *testing.T) {
			tests.SchemaTestUnionStringprefix(t, engine)
		})
	}
}
