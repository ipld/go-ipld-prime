package schema

import (
	"github.com/ipld/go-ipld-prime"
)

type TypeUnion struct {
	ts      *TypeSystem
	name    TypeName
	members []TypeName // all of these are TypeName because we ruled by fiat that unions are not allowed to use anon types (for sheer syntactic complexity boundary reasons).
	rstrat  UnionRepresentation
}

type UnionRepresentation interface{ _UnionRepresentation() }

func (UnionRepresentation_Keyed) _UnionRepresentation()        {}
func (UnionRepresentation_Kinded) _UnionRepresentation()       {}
func (UnionRepresentation_Envelope) _UnionRepresentation()     {}
func (UnionRepresentation_Inline) _UnionRepresentation()       {}
func (UnionRepresentation_StringPrefix) _UnionRepresentation() {}
func (UnionRepresentation_BytePrefix) _UnionRepresentation()   {}

type UnionRepresentation_Keyed struct {
	ts                *TypeSystem
	discriminantTable map[string]TypeName
}
type UnionRepresentation_Kinded struct {
	ts                *TypeSystem
	discriminantTable map[ipld.Kind]TypeName
}
type UnionRepresentation_Envelope struct {
	ts                *TypeSystem
	discriminantKey   string
	contentKey        string
	discriminantTable map[string]TypeName
}
type UnionRepresentation_Inline struct {
	ts                *TypeSystem
	discriminantKey   string
	discriminantTable map[string]TypeName
}
type UnionRepresentation_StringPrefix struct {
	ts                *TypeSystem
	discriminantTable map[string]TypeName
}
type UnionRepresentation_BytePrefix struct {
	ts                *TypeSystem
	discriminantTable map[string]TypeName
}

// -- Type interface satisfaction -->

var _ Type = (*TypeUnion)(nil)

func (t *TypeUnion) TypeSystem() *TypeSystem {
	return t.ts
}

func (TypeUnion) TypeKind() TypeKind {
	return TypeKind_Union
}

func (t *TypeUnion) Name() TypeName {
	return t.name
}

func (t *TypeUnion) RepresentationBehavior() ipld.Kind {
	switch t.rstrat.(type) {
	case UnionRepresentation_Keyed:
		return ipld.Kind_Map
	case UnionRepresentation_Kinded:
		return ipld.Kind_Invalid // you can't know with this one, until you see the value (and thus can see its inhabitant's behavior)!
	case UnionRepresentation_Envelope:
		return ipld.Kind_Map
	case UnionRepresentation_Inline:
		return ipld.Kind_Map
	case UnionRepresentation_StringPrefix:
		return ipld.Kind_String
	case UnionRepresentation_BytePrefix:
		return ipld.Kind_Bytes
	default:
		panic("unreachable")
	}
}

// -- specific to TypeUnion -->

func (t *TypeUnion) RepresentationStrategy() UnionRepresentation {
	return t.rstrat
}

// GetDiscriminantForType looks up the discriminant key for the given type.
// It panics if the given type is not a member of this union.
func (r UnionRepresentation_Keyed) GetDiscriminantForType(t Type) string {
	if t.TypeSystem() != r.ts {
		panic("that type isn't even from the same universe!")
	}
	for k, v := range r.discriminantTable {
		if v == t.Name() {
			return k
		}
	}
	panic("that type isn't a member of this union")
}

// GetMember returns the type info for the member that would be indicated by the given kind,
// or may return nil if that kind is not mapped to a member of this union.
func (r UnionRepresentation_Kinded) GetMember(k ipld.Kind) Type {
	if tn, exists := r.discriminantTable[k]; exists {
		return r.ts.types[TypeReference(tn)]
	}
	return nil
}
