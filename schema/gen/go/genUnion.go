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
	//
	// A note about index: in all cases the index of a member type is used, we increment it by one, to avoid using zero.
	// We do this because it's desirable to reserve the zero in the 'tag' field (if we generate one) as a sentinel value
	// (see further comments in the EmitNodeAssemblerType function);
	// and since we do it in that one case, it's just as well to do it uniformly.
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct {
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			tag uint
			{{- range $i, $member := .Type.Members }}
			x{{ add $i 1 }} _{{ $member | TypeSymbol }}
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
			case {{ add $i 1 }}:
				return &n.x{{ add $i 1 }}
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
				if n.tag != {{ add $i 1 }} {
					return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
				}
				return &n.x{{ add $i 1 }}, nil
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
			case {{ add $i 1 }}:
				return memberName__{{ dot.Type | TypeSymbol }}_{{ $member.Name }}, &n.x{{ add $i 1 }}, nil
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
	return unionBuilderGenerator{
		g.AdjCfg,
		mixins.MapAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__",
		},
		g.PkgName,
		g.Type,
	}
}

type unionBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapAssemblerTraits
	PkgName string
	Type    schema.TypeUnion
}

func (unionBuilderGenerator) IsRepr() bool { return false } // hint used in some generalized templates.

func (g unionBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	emitEmitNodeBuilderType_typical(w, g.AdjCfg, g)
}
func (g unionBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	emitNodeBuilderMethods_typical(w, g.AdjCfg, g)
}
func (g unionBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// Assemblers for unions are not unlikely those for structs or maps:
	//
	// - 'w' is the "**w**ip" pointer.
	// - 'm' is the pointer to a **m**aybe which communicates our completeness to the parent if we're a child assembler.
	//     Like any other structure, a union can be nullable in the context of some enclosing object, and we'll have the usual branches for handling that in our various Assign methods.
	// - 'state' is what it says on the tin.  Unions use maState to sequence the transitions between a new assembler, the map having been started, key insertions, value insertions, and finish.
	//     Most of this is just like the way struct and map use maState.
	//     However, we also need to guard to make sure a second entry never begins; after the first, finish is the *only* valid transition.
	//     In structs, this is done using the "set" bitfield; in maps, the state resides in the wip map itself.
	//     Unions are more like the latter: depending on which memory layout we're using, either the `na.w.tag` value, or, a non-nil `na.w.x`, is indicative that one key has been entered.
	//     (The zero value for `na.w.tag` is reserved, and all  for this reason.
	// - There is no additional state need to store "focus" (in contrast to structs);
	//     information during the AssembleValue phase about which member is selected is also just handled in `na.w.tag`, or, in the type info of `na.w.x`, again depending on memory layout strategy.
	//
	// - 'cm' is **c**hild **m**aybe and is used for the completion message from children.
	// - 'ca*' fields embed **c**hild **a**ssemblers -- these are embedded so we can yield pointers to them during recusion into child value assembly without causing new allocations.
	//     In unions, only one of these will every be used!  However, we don't know *which one* in advance, so, we have to embed them all.
	//     (It's ironic to note that if the golang compiler had an understanding of unions itself (either tagged or untagged would suffice), we could compile this down into *much* more minimal amounts of resident memory reservation.  Alas!)
	//     We elide all of the 'ca*' embeds, and instead allocate on demand, for unions with memlayout=interface mode.  (Arguably, this is overloading that config; PRs for more granular configurability welcome.)
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe
			state maState

			cm schema.Maybe
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			{{- range $i, $member := .Type.Members }}
			ca{{ add $i 1 }} _{{ $member | TypeSymbol }}__Assembler
			{{end -}}
			{{end -}}
		}

		func (na *_{{ .Type | TypeSymbol }}__Assembler) reset() {
			na.state = maState_initial
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			{{- range $i, $member := .Type.Members }}
			na.ca{{ add $i 1 }}.reset()
			{{end -}}
			{{end -}}
		}
	`, w, g.AdjCfg, g)
}
func (g unionBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	// We currently disregard sizeHint.  It's not relevant to us.
	//  We could check it strictly and emit errors; presently, we don't.
	// This method contains a branch to support MaybeUsesPtr because new memory may need to be allocated.
	//  This allocation only happens if the 'w' ptr is nil, which means we're being used on a Maybe;
	//  otherwise, the 'w' ptr should already be set, and we fill that memory location without allocating, as usual.
	// DRY: this turns out to be textually identical to the method for structs!
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) BeginMap(int) (ipld.MapAssembler, error) {
			switch *na.m {
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			case midvalue:
				panic("invalid state: it makes no sense to 'begin' twice on the same assembler!")
			}
			*na.m = midvalue
			{{- if .Type | MaybeUsesPtr }}
			if na.w == nil {
				na.w = &_{{ .Type | TypeSymbol }}{}
			}
			{{- end}}
			return na, nil
		}
	`, w, g.AdjCfg, g)
}
func (g unionBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	// It might sound a bit odd to call a union "recursive", since it's so very trivially so (no fan-out),
	//  but it's functionally accurate: the generated method should include a branch for the 'midvalue' state.
	emitNodeAssemblerMethodAssignNull_recursive(w, g.AdjCfg, g)
}
func (g unionBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	// AssignNode goes through three phases:
	// 1. is it null?  Jump over to AssignNull (which may or may not reject it).
	// 2. is it our own type?  Handle specially -- we might be able to do efficient things.
	// 3. is it the right kind to morph into us?  Do so.
	//
	// We do not set m=midvalue in phase 3 -- it shouldn't matter unless you're trying to pull off concurrent access, which is wrong and unsafe regardless.
	//
	// DRY: this turns out to be textually identical to the method for structs!  (At least, for now.  It could/should probably be optimized to get to the point faster in phase 3.)
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNode(v ipld.Node) error {
			if v.IsNull() {
				return na.AssignNull()
			}
			if v2, ok := v.(*_{{ .Type | TypeSymbol }}); ok {
				switch *na.m {
				case schema.Maybe_Value, schema.Maybe_Null:
					panic("invalid state: cannot assign into assembler that's already finished")
				case midvalue:
					panic("invalid state: cannot assign null into an assembler that's already begun working on recursive structures!")
				}
				{{- if .Type | MaybeUsesPtr }}
				if na.w == nil {
					na.w = v2
					*na.m = schema.Maybe_Value
					return nil
				}
				{{- end}}
				*na.w = *v2
				*na.m = schema.Maybe_Value
				return nil
			}
			if v.ReprKind() != ipld.ReprKind_Map {
				return ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .Type.Name }}", MethodName: "AssignNode", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: v.ReprKind()}
			}
			itr := v.MapIterator()
			for !itr.Done() {
				k, v, err := itr.Next()
				if err != nil {
					return err
				}
				if err := na.AssembleKey().AssignNode(k); err != nil {
					return err
				}
				if err := na.AssembleValue().AssignNode(v); err != nil {
					return err
				}
			}
			return na.Finish()
		}
	`, w, g.AdjCfg, g)
}
func (g unionBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// TODO SOON
}
