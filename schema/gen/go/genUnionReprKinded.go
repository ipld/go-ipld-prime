package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

var _ TypeGenerator = &unionReprKindedGenerator{}

// Kinded union representations are quite wild: their behavior varies almost completely per inhabitant,
//  and their implementation is generally delegating directly to something else,
//   rather than having an intermediate node (like most unions do, and like the type-level view of this same value will).
//
// This also means any error values can be a little weird:
//  sometimes they'll have the union's type name, but sometimes they'll have the inhabitant's type name instead;
//  this depends on whether the error is an ErrWrongKind that was found while checking the method for appropriateness on the union's inhabitant
//  versus if the error came from the union inhabitant itself after delegation occured.

func NewUnionReprKindedGenerator(pkgName string, typ *schema.TypeUnion, adjCfg *AdjunctCfg) TypeGenerator {
	return unionReprKindedGenerator{
		unionGenerator{
			adjCfg,
			mixins.MapTraits{
				pkgName,
				string(typ.Name()),
				adjCfg.TypeSymbol(typ),
			},
			pkgName,
			typ,
		},
	}
}

type unionReprKindedGenerator struct {
	unionGenerator
}

func (g unionReprKindedGenerator) GetRepresentationNodeGen() NodeGenerator {
	return unionReprKindedReprGenerator{
		g.AdjCfg,
		g.PkgName,
		g.Type,
	}
}

type unionReprKindedReprGenerator struct {
	// Note that there's no MapTraits (or any other FooTraits) mixin in this one!
	//  This is no accident: *None* of them apply!

	AdjCfg  *AdjunctCfg
	PkgName string
	Type    *schema.TypeUnion
}

func (unionReprKindedReprGenerator) IsRepr() bool { return true } // hint used in some generalized templates.

func (g unionReprKindedReprGenerator) EmitNodeType(w io.Writer) {
	// The type is structurally the same, but will have a different set of methods.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Repr _{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = &_{{ .Type | TypeSymbol }}__Repr{}
	`, w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodReprKind(w io.Writer) {
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}__Repr) ReprKind() ipld.ReprKind {
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			switch n.tag {
			{{- range $i, $member := .Type.Members }}
			case {{ add $i 1 }}:
				return {{ $member.RepresentationBehavior | KindSymbol }}
			{{- end}}
			{{- else if (eq (.AdjCfg.UnionMemlayout .Type) "interface") }}
			switch n2 := n.x.(type) {
			{{- range $i, $member := .Type.Members }}
			case {{ $member | TypeSymbol }}:
				return {{ $member.RepresentationBehavior | KindSymbol }}
			{{- end}}
			{{- end}}
			default:
				panic("unreachable")
			}
		}
	`, w, g.AdjCfg, g)
}

// A bunch of these methods could be improved by doing a gen-time switch for whether any of the possible members are the relevant kind at all;
//  currently in the cases where there's no relevant members, we generate switch blocks that are empty except for their default...
//   which works, but is arguably a little strange.
//    I haven't checked if this dummy switch has any actual performance implications:
//     I haven't tested if this produces unconditional assembly,
//     nor if it successfully removes the access of the tag,
//     though one might imagine a sufficiently clever compiler ought to do both of those things.
//     Regardless, the gsloc is reducable.  (Slightly.  There are also bigger gains to be made elsewhere, I'm sure.)

func kindedUnionNodeMethodTemplateMunge(
	methodName string, // for error messages
	methodSig string, // output literally
	condClause string, // template condition for the member this should match on
	retClause string, // clause returning the thing (how to delegate methodsig, generally)
	appropriateKind string, // for error messages
	nopeSentinel string, // for error return paths; generally the zero value for the first return type.
	nopeSentinelOnly bool, // true if this method has no error return, just the sentinel.
) string {
	// We really could just... call the methods directly (and elide the switch entirely all the time), in the case of the "interface" implementation strategy.
	//  We don't, though, because that would deprive us of getting the union type's name in the wrong-kind errors...
	//   and in addition to that being sadface in general, it would be downright unacceptable if that behavior varied based on implementation strategy.
	//
	// This error text doesn't tell us what the member kind is.  This might read weirdly.
	//  It's possible we could try to cram a description of the inhabitant into the "TypeName" since it's stringy; but unclear if that's a good idea either.
	errorClause := `return ` + nopeSentinel
	if !nopeSentinelOnly {
		errorClause += `, ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .Type.Name }}.Repr", MethodName: "` + methodName + `", AppropriateKind: ` + appropriateKind + `, ActualKind: n.ReprKind()}`
	}
	return `
		func (n *_{{ .Type | TypeSymbol }}__Repr) ` + methodSig + ` {
			{{- if (eq (.AdjCfg.UnionMemlayout .Type) "embedAll") }}
			switch n.tag {
			{{- range $i, $member := .Type.Members }}
			` + condClause + `
			case {{ add $i 1 }}:
				return n.x{{ add $i 1 }}.Representation()` + retClause + `
			{{- end}}
			{{- end}}
			{{- else if (eq (.AdjCfg.UnionMemlayout .Type) "interface") }}
			switch n2 := n.x.(type) {
			{{- range $i, $member := .Type.Members }}
			` + condClause + `
			case {{ $member | TypeSymbol }}:
				return n2.Representation()` + retClause + `
			{{- end}}
			{{- end}}
			{{- end}}
			default:
				` + errorClause + `
			}
		}
	`
}

func (g unionReprKindedReprGenerator) EmitNodeMethodLookupByString(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`LookupByString`,
		`LookupByString(key string) (ipld.Node, error)`,
		`{{- if eq $member.RepresentationBehavior.String "map" }}`,
		`.LookupByString(key)`,
		`ipld.ReprKindSet_JustMap`,
		`nil`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodLookupByIndex(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`LookupByIndex`,
		`LookupByIndex(idx int) (ipld.Node, error)`,
		`{{- if eq $member.RepresentationBehavior.String "list" }}`,
		`.LookupByIndex(idx)`,
		`ipld.ReprKindSet_JustList`,
		`nil`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodLookupByNode(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`LookupByNode`,
		`LookupByNode(key ipld.Node) (ipld.Node, error)`,
		`{{- if or (eq $member.RepresentationBehavior.String "map") (eq $member.RepresentationBehavior.String "list") }}`,
		`.LookupByNode(key)`,
		`ipld.ReprKindSet_Recursive`,
		`nil`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodLookupBySegment(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`LookupBySegment`,
		`LookupBySegment(seg ipld.PathSegment) (ipld.Node, error)`,
		`{{- if or (eq $member.RepresentationBehavior.String "map") (eq $member.RepresentationBehavior.String "list") }}`,
		`.LookupBySegment(seg)`,
		`ipld.ReprKindSet_Recursive`,
		`nil`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodMapIterator(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`MapIterator`,
		`MapIterator() ipld.MapIterator`,
		`{{- if eq $member.RepresentationBehavior.String "map" }}`,
		`.MapIterator()`,
		`ipld.ReprKindSet_JustMap`,
		`nil`,
		true,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodListIterator(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`ListIterator`,
		`ListIterator() ipld.ListIterator`,
		`{{- if eq $member.RepresentationBehavior.String "list" }}`,
		`.ListIterator()`,
		`ipld.ReprKindSet_JustList`,
		`nil`,
		true,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`Length`,
		`Length() int`,
		`{{- if or (eq $member.RepresentationBehavior.String "map") (eq $member.RepresentationBehavior.String "list") }}`,
		`.Length()`,
		`ipld.ReprKindSet_Recursive`,
		`-1`,
		true,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodIsAbsent(w io.Writer) {
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}__Repr) IsAbsent() bool {
			return false
		}
	`, w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodIsNull(w io.Writer) {
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}__Repr) IsNull() bool {
			return false
		}
	`, w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodAsBool(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`AsBool`,
		`AsBool() (bool, error)`,
		`{{- if eq $member.RepresentationBehavior.String "bool" }}`,
		`.AsBool()`,
		`ipld.ReprKindSet_JustBool`,
		`false`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodAsInt(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`AsInt`,
		`AsInt() (int, error)`,
		`{{- if eq $member.RepresentationBehavior.String "int" }}`,
		`.AsInt()`,
		`ipld.ReprKindSet_JustInt`,
		`0`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodAsFloat(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`AsFloat`,
		`AsFloat() (float64, error)`,
		`{{- if eq $member.RepresentationBehavior.String "float" }}`,
		`.AsFloat()`,
		`ipld.ReprKindSet_JustFloat`,
		`0`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`AsString`,
		`AsString() (string, error)`,
		`{{- if eq $member.RepresentationBehavior.String "string" }}`,
		`.AsString()`,
		`ipld.ReprKindSet_JustString`,
		`""`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodAsBytes(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`AsBytes`,
		`AsBytes() ([]byte, error)`,
		`{{- if eq $member.RepresentationBehavior.String "bytes" }}`,
		`.AsBytes()`,
		`ipld.ReprKindSet_JustBytes`,
		`nil`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodAsLink(w io.Writer) {
	doTemplate(kindedUnionNodeMethodTemplateMunge(
		`AsLink`,
		`AsLink() (ipld.Link, error)`,
		`{{- if eq $member.RepresentationBehavior.String "link" }}`,
		`.AsLink()`,
		`ipld.ReprKindSet_JustLink`,
		`nil`,
		false,
	), w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodeMethodPrototype(w io.Writer) {
	emitNodeMethodPrototype_typical(w, g.AdjCfg, g)
}

func (g unionReprKindedReprGenerator) EmitNodePrototypeType(w io.Writer) {
	emitNodePrototypeType_typical(w, g.AdjCfg, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g unionReprKindedReprGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return unionReprKindedReprBuilderGenerator{
		g.AdjCfg,
		g.PkgName,
		g.Type,
	}
}

type unionReprKindedReprBuilderGenerator struct {
	AdjCfg  *AdjunctCfg
	PkgName string
	Type    *schema.TypeUnion
}

func (unionReprKindedReprBuilderGenerator) IsRepr() bool { return true } // hint used in some generalized templates.

func (g unionReprKindedReprBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	emitEmitNodeBuilderType_typical(w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	emitNodeBuilderMethods_typical(w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// Much of this is familiar: the 'w', the 'm' are all as usual.
	// Some things may look a little odd here compared to all other assemblers:
	//  we're kinda halfway between what's conventionally seen for a scalar and what's conventionally seen for a recursive.
	// There's no 'maState' or 'laState'-typed fields (which feels like a scalar) because even if we end up acting like a map or list, that state is in the relevant child assembler.
	// We don't even have a 'cm' field, because we can get away with something really funky: we can just copy our own 'm' _pointer_ into children; our doneness and their doneness is the same.
	// We never have to worry about maybeism of our children; the nullable and optional modifiers aren't possible on union members.
	//  (We *do* still have to consider null values though, as null is still a kind, and thus can be routed to one of our members!)
	// 'ca' is as it is in the type-level assembler: technically, not super necessary, except that it allows minimizing the amount of work that resetting needs to do.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprAssembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe

			{{- range $i, $member := .Type.Members }}
			ca{{ add $i 1 }} {{ if (eq (dot.AdjCfg.UnionMemlayout dot.Type) "interface") }}*{{end}}_{{ $member | TypeSymbol }}__ReprAssembler
			{{end -}}
			ca uint
		}
	`, w, g.AdjCfg, g)
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) reset() {
			na.state = maState_initial
			switch na.ca {
			case 0:
				return
			{{- range $i, $member := .Type.Members }}
			case {{ add $i 1 }}:
				na.ca{{ add $i 1 }}.reset()
			{{end -}}
			default:
				panic("unreachable")
			}
		}
	`, w, g.AdjCfg, g)
}

func kindedUnionNodeAssemblerMethodTemplateMunge(
	methodSig string, condClause string, retClause string,
) string {
	// The value pointed to by `na.m` isn't modified here, because we're sharing it with the child, who should do so.
	//  This also means that value gets checked twice -- once by us, because we need to halt if we've already been used --
	//   and also a second time by the child when we delegate to it, which, unbeknownst to it, is irrelevant.
	//   I don't see a good way to remedy this shy of making more granular (unexported!) methods.  (Might be worth it.)
	//   This probably also isn't the same for all of the assembler methods: the methods we delegate to aren't doing as many check branches when they're for scalars,
	//    because they expected to be used in contexts where many values of the 'm' enum aren't reachable -- an expectation we've suddenly subverted with this path!
	return `
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) ` + methodSig + ` {
			switch *na.m {
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			case midvalue:
				panic("invalid state: cannot assign into assembler that's already working on a larger structure!")
			}
			{{- range $i, $member := .Type.Members }}
			` + condClause + `
			{{- if dot.Type | MaybeUsesPtr }}
			if na.w == nil {
				na.w = &_{{ dot.Type | TypeSymbol }}{}
			}
			{{- end}}
			na.ca = {{ add $i 1 }}
			{{- if (eq (dot.AdjCfg.UnionMemlayout dot.Type) "embedAll") }}
			na.w.tag = {{ add $i 1 }}
			na.ca{{ add $i 1 }}.w = &na.w.x{{ add $i 1 }}
			na.ca{{ add $i 1 }}.m = na.m
			return na.ca{{ add $i 1 }}.BeginMap(sizeHint)
			{{- else if (eq (dot.AdjCfg.UnionMemlayout dot.Type) "interface") }}
			x := &_{{ $member | TypeSymbol }}{}
			na.w.x = x
			if na.ca{{ add $i 1 }} == nil {
				na.ca{{ add $i 1 }} = &_{{ $member | TypeSymbol }}__ReprAssembler{}
			}
			na.ca{{ add $i 1 }}.w = x
			na.ca{{ add $i 1 }}.m = na.m
			return na.ca{{ add $i 1 }}` + retClause + `
			{{- end}}
			{{- end}}
			{{- end}}
			// TODO i think you finally Need a method for if-no-members-match-this-kind for the default rejection to compile this time.
			return nil, ipld.ErrWrongKind{doozy}
		}
	`
}

func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`BeginMap(sizeHint int) (ipld.MapAssembler, error)`,
		`{{- if eq $member.RepresentationBehavior.String "map" }}`,
		`.BeginMap(sizeHint)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`BeginList(sizeHint int) (ipld.ListAssembler, error)`,
		`{{- if eq $member.RepresentationBehavior.String "list" }}`,
		`.BeginList(sizeHint)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	// TODO: I think this may need some special handling to account for if our union is itself used in a nullable circumstance; that should overrule this behavior.
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`AssignNull() error `,
		`{{- if eq $member.RepresentationBehavior.String "null" }}`,
		`.AssignNull()`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`AssignBool(v bool) error `,
		`{{- if eq $member.RepresentationBehavior.String "bool" }}`,
		`.AssignBool(v)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`AssignInt(v int) error `,
		`{{- if eq $member.RepresentationBehavior.String "int" }}`,
		`.AssignInt(v)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`AssignFloat(v float64) error `,
		`{{- if eq $member.RepresentationBehavior.String "float" }}`,
		`.AssignFloat(v)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`AssignString(v string) error `,
		`{{- if eq $member.RepresentationBehavior.String "string" }}`,
		`.AssignString(v)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`AssignBytes(v []byte) error `,
		`{{- if eq $member.RepresentationBehavior.String "bytes" }}`,
		`.AssignBytes(v)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	doTemplate(kindedUnionNodeAssemblerMethodTemplateMunge(
		`AssignLink(v ipld.Link) error `,
		`{{- if eq $member.RepresentationBehavior.String "link" }}`,
		`.AssignLink(v)`,
	), w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	// TODO this is too wild for me at the moment, come back to it shortly
	// it's basically got some of the body of kindedUnionNodeAssemblerMethodTemplateMunge, but repeated many more times.
	// it also needs to handle nulls gingerly.
	// and also handle pumping the full copy in the case of lists or maps.

	// this is gonna have a fun ErrWrongKind value too -- we might actually have to make a non-static set of acceptable kinds :D  that's a first.
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerMethodPrototype(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) Prototype() ipld.NodePrototype {
			return _{{ .Type | TypeSymbol }}__ReprPrototype{}
		}
	`, w, g.AdjCfg, g)
}
func (g unionReprKindedReprBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// somewhat shockingly: nothing.
}
