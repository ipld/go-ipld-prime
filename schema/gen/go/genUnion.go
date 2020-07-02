package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

// The generator for unions is a bit more wild than most others:
// it has at three major branches for how its internals are laid out:
//
//   - all possible children are embedded.
//   - all possible children are pointers... in which case we collapse to one interface resident.
//       (n.b. this does give up some inlining potential as well as gives up on alloc amortization, but it does make resident memory size minimal.)
//   - some children are emebedded and some are pointers, and of the latter set, they may be either in one interface field or several discrete pointers.
//       (discrete fields of pointer type makes inlining possible in some paths, whereas an interface field blocks it).
//
// ... We're not doing that last one at all right now.  The pareto-prevalence of these concerns is extremely low compared to the effort required.
// But the first two are both very reasonable, and both are often wanted.
//
// These choices are made from adjunct config (which should make sense, because they're clearly all "golang" details -- not type semantics).
// We still tackle all the generation for all these strategies this in one file,
//  because all of the interfaces we export are the same, regardless of the internals (and it just seems easiest to do this way).

type unionGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapTraits
	PkgName string
	Type    schema.TypeUnion
}

func (unionGenerator) IsRepr() bool { return false } // hint used in some generalized templates.

// --- native content and specializations --->

func (g unionGenerator) EmitNativeType(w io.Writer) {
	// We generate *two* types: a struct which acts as the union node,
	// and also an interface which covers the members (and has an unexported marker function to make sure the set can't be extended).
	//
	// The interface *mostly* isn't used... except for in the return type of a speciated function which can be used to do golang-native type switches.
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct {
			{{- with (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			{{- range $member := .Type.Members }}
			x_{{ $member.Name }} _{{ $member | TypeSymbol }}
			{{- end}}
			{{- end}}
			{{- with (eq (.AdjCfg.UnionMemlayout .Type) "interface") }}
			x _{{ .Type | TypeSymbol }}__iface
			{{- end}}
		}
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}

		type _{{ .Type | TypeSymbol }}__iface interface {
			_{{ .Type | TypeSymbol }}__member()
		}

		{{- range $member := .Type.Members }}
		func (_{{ $member | TypeSymbol }}) _{{ .dot.Type | TypeSymbol }}__member() {}
		{{- end}}
	`, w, g.AdjCfg, g)
}
