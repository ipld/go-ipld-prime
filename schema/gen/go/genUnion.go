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
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct {
			{{ with (eq .AdjCfg.UnionMemlayout "embedAll") -}}

			{{- end}}
			{{ with (eq .AdjCfg.UnionMemlayout "interface") -}}

			{{- end}}
		}
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}


		{{ with (eq .AdjCfg.UnionMemlayout "interface") -}}
			// ... is there any utility to making an interface type?
			// internally: no.
			// for export and use: unclear.
			//   i don't think we need to kowtow to burntsushi's sumtype checker.  we could (and it would be pleasing to do so); but we can make our own just as well.
			//   how would you use it?  would there be a method for unboxing the structptr into the interface type?  and vice versa (which would incur an alloc)?
			//     if assignNode was going to work, ... actually it could just take the member type "concretely", they already implement Node, and we could just figure it out.  That'd be fine.
			//   is there... any reason a programmer would prefer to use the interface, though?
			//     if they wanted to make their own typeswitch, maybe?
			//       is that going to be more efficient that doing a switch using an enum we generate, and then calling a function that does a cast explicitly?  probably.  heck.
			//         do we need to have an *exported* interface for this to work, though?  I think not.  so that might be a solid choice.
		{{- end}}
	`, w, g.AdjCfg, g)
}
