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
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			tag uint
			{{- range $i, $member := .Type.Members }}
			x{{ $i }} _{{ $member | TypeSymbol }}
			{{- end}}
			{{- else if (eq (.AdjCfg.UnionMemlayout .Type) "interface") }}
			x _{{ .Type | TypeSymbol }}__iface
			{{- end}}
		}
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}

		type _{{ .Type | TypeSymbol }}__iface interface {
			_{{ .Type | TypeSymbol }}__member()
		}

		{{- range $member := .Type.Members }}
		func (_{{ $member | TypeSymbol }}) _{{ dot.Type | TypeSymbol }}__member() {}
		{{- end}}
	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNativeAccessors(w io.Writer) {
	doTemplate(`
		func (n _{{ .Type | TypeSymbol }}) AsInterface() _{{ .Type | TypeSymbol }}__iface {
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			switch n.tag {
			{{- range $i, $member := .Type.Members }}
			case {{ $i }}:
				return &n.x{{ $i }}
			{{- end}}
			default:
				panic("invalid union state; how did you create this object?")
			}
			{{- else if (eq (.AdjCfg.UnionMemlayout .Type) "interface") }}
			return n.x
			{{- end}}
		}
	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNativeBuilder(w io.Writer) {
	// Unclear as yet what should go here.
}

func (g unionGenerator) EmitNativeMaybe(w io.Writer) {
	emitNativeMaybe(w, g.AdjCfg, g)
}

// --- type info --->

func (g unionGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode interface satisfaction --->

func (g unionGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	emitTypicalTypedNodeMethodRepresentation(w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g unionGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.

	// We do, however, want some constants for our member names;
	//  they'll make iterators able to work faster.  So let's emit those.
	// These are a bit perplexing, because they're... type names.
	//  However, oddly enough, we don't have type names available *as nodes* anywhere else centrally available,
	//   so... we generate some values for them here with scoped identifers and get on with it.
	//    Maybe this could be elided with future work.
	doTemplate(`
		var (
			{{- range $member := .Type.Members }}
			memberName__{{ dot.Type | TypeSymbol }}_{{ $member.Name }} = _String{"{{ $member.Name }}"}
			{{- end }}
		)
	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNodeTypeAssertions(w io.Writer) {
	emitNodeTypeAssertions_typical(w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNodeMethodLookupByString(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupByString(key string) (ipld.Node, error) {
			switch key {
			{{- range $i, $member := .Type.Members }}
			case "{{ $member.Name }}":
				{{- if (eq (dot.AdjCfg.UnionMemlayout dot.Type) "embedAll") }}
				if n.tag != {{ $i }} {
					return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
				}
				return &n.x{{ $i }}, nil
				{{- else if (eq (dot.AdjCfg.UnionMemlayout dot.Type) "interface") }}
				if _, ok := n.x.({{ $member | TypeSymbol }}); !ok {
					return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
				}
				return n.x, nil
				{{- end}}
			{{- end}}
			default:
				return nil, schema.ErrNoSuchField{Type: nil /*TODO*/, FieldName: key}
			}
		}
	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNodeMethodLookupByNode(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupByNode(key ipld.Node) (ipld.Node, error) {
			ks, err := key.AsString()
			if err != nil {
				return nil, err
			}
			return n.LookupByString(ks)
		}
	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNodeMethodMapIterator(w io.Writer) {
	// This is kind of a hilarious "iterator": it has to count all the way up to... 1.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) MapIterator() ipld.MapIterator {
			return &_{{ .Type | TypeSymbol }}__MapItr{n, false}
		}

		type _{{ .Type | TypeSymbol }}__MapItr struct {
			n {{ .Type | TypeSymbol }}
			done bool
		}

		func (itr *_{{ .Type | TypeSymbol }}__MapItr) Next() (k ipld.Node, v ipld.Node, _ error) {
			if itr.done {
				return nil, nil, ipld.ErrIteratorOverread{}
			}
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			switch itr.n.tag {
			{{- range $i, $member := .Type.Members }}
			case {{ $i }}:
				return memberName__{{ dot.Type | TypeSymbol }}_{{ $member.Name }}, &n.x{{ $i }}, nil
			{{- end}}
			{{- else if (eq (.AdjCfg.UnionMemlayout .Type) "interface") }}
			switch itr.n.x.(type) {
			{{- range $member := .Type.Members }}
			case {{ $member | TypeSymbol }}:
				return memberName__{{ dot.Type | TypeSymbol }}_{{ $member.Name }}, n.x, nil
			{{- end}}
			{{- end}}
			default:
				panic("unreachable")
			}
			itr.done = true
			return
		}
		func (itr *_{{ .Type | TypeSymbol }}__MapItr) Done() bool {
			return itr.done
		}

	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Length() int {
			return 1
		}
	`, w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNodeMethodPrototype(w io.Writer) {
	emitNodeMethodPrototype_typical(w, g.AdjCfg, g)
}

func (g unionGenerator) EmitNodePrototypeType(w io.Writer) {
	emitNodePrototypeType_typical(w, g.AdjCfg, g)
}

func (g unionGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return nil /* TODO */
}
