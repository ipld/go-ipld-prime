package gengo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestUnionKeyed(t *testing.T) {
	t.Parallel()

	for _, engine := range []*genAndCompileEngine{
		{
			subtestName: "union-using-embed",
			prefix:      "union-keyed-using-embed",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"StrStr": "embedAll"},
			},
		},
		{
			subtestName: "union-using-ptr",
			prefix:      "union-keyed-using-interface",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"StrStr": "interface"},
			},
		},
	} {
		t.Run(engine.subtestName, func(t *testing.T) {
			tests.SchemaTestUnionKeyed(t, engine)
		})
	}
}

func TestUnionKeyedComplexChildren(t *testing.T) {
	t.Parallel()

	for _, engine := range []*genAndCompileEngine{
		{
			subtestName: "union-using-embed",
			prefix:      "union-keyed-complex-child-using-embed",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"WheeUnion": "embedAll"},
			},
		},
		{
			subtestName: "union-using-interface",
			prefix:      "union-keyed-complex-child-using-interface",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"WheeUnion": "interface"},
			},
		},
	} {
		t.Run(engine.subtestName, func(t *testing.T) {
			tests.SchemaTestUnionKeyedComplexChildren(t, engine)
		})
	}
}

func TestUnionKeyedReset(t *testing.T) {
	t.Parallel()

	for _, engine := range []*genAndCompileEngine{
		{
			subtestName: "union-using-embed",
			prefix:      "union-keyed-reset-using-embed",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"WheeUnion": "embedAll"},
			},
		},
		{
			subtestName: "union-using-interface",
			prefix:      "union-keyed-reset-using-interface",
			adjCfg: AdjunctCfg{
				CfgUnionMemlayout: map[schema.TypeName]string{"WheeUnion": "interface"},
			},
		},
	} {
		t.Run(engine.subtestName, func(t *testing.T) {
			tests.SchemaTestUnionKeyedReset(t, engine)
		})
	}
}
