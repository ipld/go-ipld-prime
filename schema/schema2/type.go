package schema

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

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
// This is a closed union; you can switch upon the above members without
// including a default case.  The membership is closed by the unexported
// '_Type' method; you may use the BurntSushi/go-sumtype tool to check
// your switches for completeness.
//
// Many interesting properties of each Type are only defined for that specific
// type, so it's typical to use a type switch to handle each type of Type.
// (Your humble author is truly sorry for the word-mash that results from
// attempting to describe the types that describe the typesystem.Type.)
//
// For example, to inspect the kind of fields in a struct: you might
// cast a `Type` interface into `TypeStruct`, and then the `Fields()` on
// that `TypeStruct` can be inspected.  (`Fields()` isn't defined for any
// other kind of Type.)
type Type interface {
	// Unexported marker method to force the union closed.
	_Type()

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
	// Note that a schema.TypeKind is a different enum than datamodel.Kind;
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
	RepresentationBehavior() datamodel.Kind
}
