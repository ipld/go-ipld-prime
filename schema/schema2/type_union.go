package schema

import (
	"github.com/ipld/go-ipld-prime"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeUnion struct {
	name TypeName
	dmt  schemadmt.TypeUnion
	ts   *TypeSystem
}

type UnionRepresentation interface{ _UnionRepresentation() }

func (UnionRepresentation_Keyed) _UnionRepresentation()        {}
func (UnionRepresentation_Kinded) _UnionRepresentation()       {}
func (UnionRepresentation_Envelope) _UnionRepresentation()     {}
func (UnionRepresentation_Inline) _UnionRepresentation()       {}
func (UnionRepresentation_StringPrefix) _UnionRepresentation() {}
func (UnionRepresentation_BytePrefix) _UnionRepresentation()   {}

type UnionRepresentation_Keyed struct {
	ts  *TypeSystem
	dmt schemadmt.UnionRepresentation_Keyed
}
type UnionRepresentation_Kinded struct {
	ts  *TypeSystem
	dmt schemadmt.UnionRepresentation_Kinded
}
type UnionRepresentation_Envelope struct {
	ts  *TypeSystem
	dmt schemadmt.UnionRepresentation_Envelope
}
type UnionRepresentation_Inline struct {
	ts  *TypeSystem
	dmt schemadmt.UnionRepresentation_Inline
}
type UnionRepresentation_StringPrefix struct {
	ts  *TypeSystem
	dmt schemadmt.UnionRepresentation_StringPrefix
}
type UnionRepresentation_BytePrefix struct {
	ts  *TypeSystem
	dmt schemadmt.UnionRepresentation_BytePrefix
}

// -- schema.Type interface satisfaction -->

var _ Type = (*TypeUnion)(nil)

func (t *TypeUnion) _Type() {}

func (t *TypeUnion) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeUnion) Kind() Kind {
	return Kind_Union
}

func (t *TypeUnion) Name() TypeName {
	return t.name
}

func (t TypeUnion) RepresentationBehavior() ipld.ReprKind {
	switch t.dmt.FieldRepresentation().AsInterface().(type) {
	case schemadmt.UnionRepresentation_Keyed:
		return ipld.ReprKind_Map
	case schemadmt.UnionRepresentation_Kinded:
		return ipld.ReprKind_Invalid // you can't know with this one, until you see the value (and thus can see its inhabitant's behavior)!
	case schemadmt.UnionRepresentation_Envelope:
		return ipld.ReprKind_Map
	case schemadmt.UnionRepresentation_Inline:
		return ipld.ReprKind_Map
	case schemadmt.UnionRepresentation_StringPrefix:
		return ipld.ReprKind_String
	case schemadmt.UnionRepresentation_BytePrefix:
		return ipld.ReprKind_Bytes
	default:
		panic("unreachable")
	}
}

// -- specific to TypeUnion -->

func (t *TypeUnion) RepresentationStrategy() UnionRepresentation {
	switch x := t.dmt.FieldRepresentation().AsInterface().(type) {
	case schemadmt.UnionRepresentation_Keyed:
		return UnionRepresentation_Keyed{t.ts, x}
	case schemadmt.UnionRepresentation_Kinded:
		return UnionRepresentation_Kinded{t.ts, x}
	case schemadmt.UnionRepresentation_Envelope:
		return UnionRepresentation_Envelope{t.ts, x}
	case schemadmt.UnionRepresentation_Inline:
		return UnionRepresentation_Inline{t.ts, x}
	case schemadmt.UnionRepresentation_StringPrefix:
		return UnionRepresentation_StringPrefix{t.ts, x}
	case schemadmt.UnionRepresentation_BytePrefix:
		return UnionRepresentation_BytePrefix{t.ts, x}
	default:
		panic("unreachable")
	}
}

// GetDiscriminantForType looks up the descriminant key for the given type.
// It panics if the given type is not a member of this union.
func (r UnionRepresentation_Keyed) GetDiscriminantForType(t Type) string {
	if t.TypeSystem() != r.ts {
		panic("that type isn't even from the same universe!")
	}
	for itr := r.dmt.Iterator(); !itr.Done(); {
		k, v := itr.Next()
		if v == t.Name() {
			return k.String()
		}
	}
	panic("that type isn't a member of this union")
}

// GetMember returns type info for the member matching the kind argument,
// or may return nil if that kind is not mapped to a member of this union.
func (r UnionRepresentation_Kinded) GetMember(k ipld.ReprKind) Type {
	rkdmt, _ := schemadmt.Type.RepresentationKind.FromString(k.String()) // FUTURE: this is currently awkward because we used a string where we should use an enum; this can be fixed when codegen for enums is implemented.
	tn := r.dmt.Lookup(rkdmt)
	if tn == nil {
		return nil
	}
	return r.ts.types[tn.TypeReference()]
}
