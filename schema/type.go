package schema

import (
	ipld "github.com/ipld/go-ipld-prime"
)

// TypeName is a string that names a type.
// TypeName is restricted to UTF-8 numbers and letter, must start with a letter,
// excludes whitespace, and excludes multiple consecutive underscores.
//
// More specifically, the definitions used in https://golang.org/ref/spec#Identifiers
// apply for defining numbers and letters.  We don't recommend pushing the limits
// and corner cases in schemas you author, either, however; tooling in other
// langauges may be made more difficult to use if you do so.
type TypeName string

func (tn TypeName) String() string { return string(tn) }

// TypeReference is a string that's either a TypeName or a computed string from an InlineDefn.
// This string is often useful as a map key.
//
// The computed string for an InlineDefn happens to match the IPLD Schema DSL syntax,
// but it would be very odd for any code to depend on that detail.
type TypeReference string

// typesystem.Type is an union interface; each of the `Type*` concrete types
// in this package are one of its members.
//
// Specifically,
//
// 	TypeBool
// 	TypeString
// 	TypeBytes
// 	TypeInt
// 	TypeFloat
// 	TypeMap
// 	TypeList
// 	TypeLink
// 	TypeUnion
// 	TypeStruct
// 	TypeEnum
//
// are all of the kinds of Type.
//
// Many typekinds have interesting properties which only defined for that specific typekind,
// For example, map typekinds have additional type info about their keys and values,
// while struct types have additional information about their fields.
// Since none of these are uniformly true of all types, they aren't in this interface,
// and it's typical to use a type switch to refine to one of the more specific "Type*"
// interfaces to get those more specific pieces of information
//
// For example, to inspect the kind of fields in a struct: you might
// cast a `Type` interface into `TypeStruct`, and then the `Fields()` on
// that `TypeStruct` can be inspected.  (`Fields()` isn't defined for any
// other kind of Type.)
type Type interface {
	// Returns a pointer to the TypeSystem this Type is a member of.
	TypeSystem() *TypeSystem

	// Returns the string name of the Type.  This name is unique within the
	// universe this type is a member of, *unless* this type is Anonymous,
	// in which case a string describing the type will still be returned, but
	// that string will not be required to be unique.
	Name() TypeName

	// Returns the TypeKind of this Type.
	//
	// The returned value is a 1:1 association with which of the concrete
	// "schema.Type*" structs this interface can be cast to.
	//
	// Note that a schema.TypeKind is a different enum than ipld.Kind;
	// and furthermore, there's no strict relationship between them.
	// schema.TypedNode values can be described by *two* distinct Kinds:
	// one which describes how the Node itself will act,
	// and another which describes how the Node presents for serialization.
	// For some combinations of Type and representation strategy, one or both
	// of the Kinds can be determined statically; but not always:
	// it can sometimes be necessary to inspect the value quite concretely
	// (e.g., `schema.TypedNode{}.Representation().Kind()`) in order to find
	// out exactly how a node will be serialized!  This is because some types
	// can vary in representation kind based on their value (specifically,
	// kinded-representation unions have this property).
	TypeKind() TypeKind

	// RepresentationBehavior returns a description of how the representation
	// of this type will behave in terms of the IPLD Data Model.
	// This property varies based on the representation strategy of a type.
	//
	// In one case, the representation behavior cannot be known statically,
	// and varies based on the data: kinded unions have this trait.
	//
	// This property is used by kinded unions, which require that their members
	// all have distinct representation behavior.
	// (It follows that a kinded union cannot have another kinded union as a member.)
	//
	// You may also be interested in a related property that might have been called "TypeBehavior".
	// However, this method doesn't exist, because it's a deterministic property of `TypeKind()`!
	// You can use `TypeKind.ActsLike()` to get type-level behavioral information.
	RepresentationBehavior() ipld.Kind
}
