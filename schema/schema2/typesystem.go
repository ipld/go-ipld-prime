package schema

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
	schemadmt "github.com/ipld/go-ipld-prime/schema/dmt"
)

type TypeSystem struct {
	// Mind the key type here: TypeReference, not TypeName.
	// The key might be a computed anon "name" which is not actually a valid type name itself.
	types map[TypeReference]Type

	// TODO: we should probably have iterable orders ready, both for named types and all types.
	//  We can derive these afresh from the dmt every time, but we either should have an exported method for that, or even just compute it eagerly and cache it.
}

func BuildTypeSystem(schdmt schemadmt.Schema) (*TypeSystem, []error) {
	ts := &TypeSystem{
		types: make(map[TypeReference]Type, schdmt.FieldTypes().Length()),
	}
	var ee []error

	// Iterate over all the types, creating the reified forms of them as we go.
	//  Some forms of validation of the data are already done by nature of the schema; others will require more work here.
	//   In general: rules which stretch across multiple types (especially, if they're graph properties) can't be implemented in schemas alone, and so end up here.
	//  Any type that has some kind of recursion (maps, lists, structs, unions, links with target type info) causes a lookup to see if the referenced types exist.
	//   FUTURE: we'll need most recursives to additionally do some checks for correct composition of nullability after the introduction of unit types.
	//  We manage to avoid the need for a two-pass system because we just give all reified types a pointer to the typesystem aggregate;
	//   this means each of them still just stores the type names, and if asked for another type pointer, looks it up on the fly (by which time it's available).
	//  Some kinds of types involve other more specific checks,
	//   such as maps verifying that their keys are stringable (which is a rule we enforce for reasons relating to pathing),
	//   and unions verifying that all their discriminant tables are complete (which is a rule that's necessary for sanity!),
	//   and etc.
	typesdmt := schdmt.FieldTypes()
	for itr := typesdmt.Iterator(); !itr.Done(); {
		tn, t := itr.Next()
		switch t2 := t.AsInterface().(type) {
		case schemadmt.TypeEnum:
			// Verify that:
			// - each value in the enumeration has an entry in its representation table.
			// - each of the representation values is distinct.  Enum representation tables are keyed by the enum value, so we have to check value uniqueness.
			if t2.FieldRepresentation().Length() != t2.FieldMembers().Length() {
				ee = append(ee, fmt.Errorf("type %s representation details must contain exactly one discriminant for each member value", tn))
				continue
			}
			switch r := t2.FieldRepresentation().AsInterface().(type) {
			case schemadmt.EnumRepresentation_String:
				vs := map[string]struct{}{}
				for itr := r.Iterator(); !itr.Done(); {
					k, v := itr.Next()
					if t2.FieldMembers().Lookup(k) == nil {
						ee = append(ee, fmt.Errorf("type %s representation contains info talking about a %q member value but there's no such member", tn, k))
					}
					if _, exists := vs[v.String()]; exists {
						ee = append(ee, fmt.Errorf("type %s representation contains a discriminant (%q) more than once", tn, v.String()))
					}
					vs[v.String()] = struct{}{}
				}
			case schemadmt.EnumRepresentation_Int:
				vs := map[int64]struct{}{}
				for itr := r.Iterator(); !itr.Done(); {
					k, v := itr.Next()
					if t2.FieldMembers().Lookup(k) == nil {
						ee = append(ee, fmt.Errorf("type %s representation contains info talking about a %q member value but there's no such member", tn, k))
					}
					if _, exists := vs[v.Int()]; exists {
						ee = append(ee, fmt.Errorf("type %s representation contains a discriminant (%q) more than once", tn, v.Int()))
					}
					vs[v.Int()] = struct{}{}
				}
			}
		case schemadmt.TypeCopy:
			panic("no support for 'copy' types.  I might want to reneg on whether these are even part of the schema dmt.")
		default:
			panic("unreachable")
		}
	}

	// Only return the assembled TypeSystem value if we encountered no errors.
	//  If we encountered errors, the TypeSystem is partially constructed and many of its contents cannot uphold their contracts, so it's better not to expose it.
	if ee == nil {
		return ts, nil
	}
	return nil, ee
}
