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
		case schemadmt.TypeUnion:
			// Verify... oh boy.  A lot of things; and each representation strategy has different constraints:
			// - for everyone: that all referenced member types exist.
			// - for everyone (but in distinctive ways): that each member type has a discriminant!
			// - for keyed unions: that's sufficient (discriminant uniqueness already enforced by map).
			// - for kinded unions: validate that that the stated kind actually matches what each type's representation kind is.
			//   - surprisingly, unions can nest... as long as they're not kinded.  (In theory, kinded union nesting could be defined, as long as their discriminants are non-overlapping, but... why would you want this?)
			// - for envelope unions: quick sanity check that discriminantKey and contentKey are distinct.
			// - for inline unions: validate that all members have map kinds...
			//   - and more specifically are map or struct types (other unions are not allowed because they wouldn't fit together validly anyway)
			//   - and for structs, validate in advance that the discriminant key doesn't collide with any field names in any of the structs.
			// - for stringprefix unions: that's sufficient (discriminant uniqueness already enforced by map).
			// - for byteprefix unions: ... we'll come back to this later.

			// Check for member type reference existence first.
			//  Build up a spare list of those type names in the process; we'll scratch stuff back off of it in a moment.
			members := make([]schemadmt.TypeName, 0, t2.FieldMembers().Length())
			missingTypes := false
			for itr := t2.FieldMembers().Iterator(); !itr.Done(); {
				_, tndmt := itr.Next()
				mtdmt := typesdmt.Lookup(tndmt)
				if mtdmt == nil {
					missingTypes = true
					ee = append(ee, fmt.Errorf("type %s refers to missing type %s as a member", tn, tndmt))
				}
				members = append(members, tndmt)
			}
			// Skip the rest of the checks if there were any missing type references.
			//  We'll be inspecting the value types more deeply and it's simpler to do that work while presuming everything is at least defined.
			if missingTypes {
				continue
			}

			// Do the per-representation-strategy checks.
			//  Every representation strategy a check that there's a discriminant for every member (though they require slightly different setup).
			//  Some representation strategies also include quite a few more checks.
			switch r := t2.FieldRepresentation().AsInterface().(type) {
			case schemadmt.UnionRepresentation_Keyed:
				checkUnionDiscriminantInfo(tn, members, r, &ee)
			case schemadmt.UnionRepresentation_Kinded:
				checkUnionDiscriminantInfo(tn, members, r, &ee)
				for itr := r.Iterator(); !itr.Done(); {
					k, v := itr.Next()
					// In the switch ahead, we briefly create the reified type for each member, just so we can use that to ask it its representation.
					//  We then let that data fall right back into the abyss.  The compiler should inline and optimize all this reasonably well.
					//  We create these temporary things rather than looking in the typesystem map we're accumulating because it makes the process work correctly regardless of order.
					//  For some of the kinds, this is fairly overkill (we know that the representation behavior of a bool type is bool because it doesn't have any other representation strategies!)
					//   but I've ground the whole thing out in a consistent way anyway.
					var mkind ipld.Kind
					switch t3 := typesdmt.Lookup(v).AsInterface().(type) {
					case schemadmt.TypeBool:
						mkind = TypeBool{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeString:
						mkind = TypeString{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeBytes:
						mkind = TypeBytes{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeInt:
						mkind = TypeInt{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeFloat:
						mkind = TypeFloat{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeMap:
						mkind = TypeMap{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeList:
						mkind = TypeList{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeLink:
						mkind = TypeLink{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeUnion:
						mkind = TypeUnion{dmt: t3}.RepresentationBehavior() // this actually flies!  it will yield Kind_Invalid for a kinded union, though, which we'll treat with a special error message.
					case schemadmt.TypeStruct:
						mkind = TypeStruct{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeEnum:
						mkind = TypeEnum{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeCopy:
						panic("no support for 'copy' types.  I might want to reneg on whether these are even part of the schema dmt.")
					default:
						panic("unreachable")
					}
					// TODO RepresentationKind is supposed to be an enum, but is not presently generated as such.  This block's use of `k` as a string should turn into something cleaner when enum gen is implemented and used for RepresentationKind.
					if mkind == ipld.Kind_Invalid {
						ee = append(ee, fmt.Errorf("kinded union %s declares a %s kind should be received as type %s, which is not sensible because that type is also a kinded union", tn, k, v))
					} else if k.String() != mkind.String() {
						ee = append(ee, fmt.Errorf("kinded union %s declares a %s kind should be received as type %s, but that type's representation kind is %s", tn, k, v, mkind))
					}
				}
			case schemadmt.UnionRepresentation_Envelope:
				checkUnionDiscriminantInfo(tn, members, r.FieldDiscriminantTable(), &ee)
				if r.FieldContentKey().String() == r.FieldDiscriminantKey().String() {
					ee = append(ee, fmt.Errorf("union %s has representation strategy envelope with conflicting content key and discriminant key", tn))
				}
			case schemadmt.UnionRepresentation_Inline:
				checkUnionDiscriminantInfo(tn, members, r.FieldDiscriminantTable(), &ee)
				for itr := r.FieldDiscriminantTable().Iterator(); !itr.Done(); {
					_, v := itr.Next()
					// As with the switch above which handles kinded union members, we go for the full destructuring here.
					//  It's slightly overkill considering that most of the type kinds will flatly error in practice, but consistency is nice.
					var mkind ipld.Kind
					switch t3 := typesdmt.Lookup(v).AsInterface().(type) {
					case schemadmt.TypeBool:
						mkind = TypeBool{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeString:
						mkind = TypeString{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeBytes:
						mkind = TypeBytes{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeInt:
						mkind = TypeInt{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeFloat:
						mkind = TypeFloat{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeMap:
						// For maps, we check the representation strategy -- it still has to be mappy! -- but that's it.
						//  Unlike for structs, where we can check ahead of time for field names which would collide with the discriminant key, with maps we're just stuck with that being a problem we can only discover at runtime.
						mkind = TypeMap{dmt: t3}.RepresentationBehavior()
					case schemadmt.TypeList:
						mkind = TypeList{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeLink:
						mkind = TypeLink{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeUnion:
						ee = append(ee, fmt.Errorf("union %s has representation strategy inline, which can't sensibly compose with any other union strategy, so %s (which is another union type) is not a valid member", tn, v))
						continue // kindcheck doesn't actually matter in this case; the error here isn't conditional on that.
					case schemadmt.TypeStruct:
						// Check representation strategy first.  Still has to be mappy.
						t4 := TypeStruct{dmt: t3}
						if t4.RepresentationBehavior() != ipld.Kind_Map {
							goto kindcheck // it'll fail, of course, but this goto DRY's the error message.
						}

						// Check for field name collisions.
						//  This uses the (temporarily) reified struct type info, so we can reuse that code which deals with rename directives.
						switch r2 := t4.RepresentationStrategy().(type) {
						case StructRepresentation_Map:
							for _, f := range t4.Fields() {
								if r.FieldDiscriminantKey().String() == r2.GetFieldKey(f) {
									ee = append(ee, fmt.Errorf("union %s has representation strategy inline, and %s is not a valid member for it because it has a field that collides with discriminantKey when represented", tn, v))
								}
							}
						default:
							panic("unreachable") // We know that the none of the other struct representation strategies result in a map kind.
						}

						continue // kindcheck already done in a unique way in this case.
					case schemadmt.TypeEnum:
						mkind = TypeEnum{dmt: t3}.RepresentationBehavior()
						goto kindcheck
					case schemadmt.TypeCopy:
						panic("no support for 'copy' types.  I might want to reneg on whether these are even part of the schema dmt.")
					default:
						panic("unreachable")
					}
				kindcheck:
					if mkind != ipld.Kind_Map {
						ee = append(ee, fmt.Errorf("union %s has representation strategy inline, which requires all members have map representations, so %s (which has representation kind %s) is not a valid member", tn, v, mkind))
					}
				}
			case schemadmt.UnionRepresentation_StringPrefix:
				checkUnionDiscriminantInfo(tn, members, r, &ee)
			case schemadmt.UnionRepresentation_BytePrefix:
				panic("nyi") // TODO byteprefix needs spec work.
			}
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

// checkUnionDiscriminantInfo verifies that every member in the list
// appears exactly once as a value in the discriminants map, and nothing else appears in the map.
// Errors are appended to ee.
// The members slice is destructively mutated.
// The typename parameter is purely for the use in error messages.
//
// The discriminantsMap is an untyped Node because it turns out convenient to do that way:
// we happen to know all the different union representations have a map *somewhere* for this,
// but its position and key types vary.  Untyped access lets us write more reusable code in this case.
func checkUnionDiscriminantInfo(tn TypeName, members []schemadmt.TypeName, discriminantsMap ipld.Node, ee *[]error) {
	for itr := discriminantsMap.MapIterator(); !itr.Done(); {
		_, v, _ := itr.Next()
		found := false
		for i, v2 := range members {
			if v == v2 {
				if found {
					*ee = append(*ee, fmt.Errorf("type %s representation details has more than one discriminant pointing to member type %s", tn, v2))
				}
				found = true
				members[i] = nil
			}
		}
		if !found {
			*ee = append(*ee, fmt.Errorf("type %s representation details include a discriminant refering to a non-member type %s", tn, v))
		}
	}
	for _, m := range members {
		if m != nil {
			*ee = append(*ee, fmt.Errorf("type %s representation details is missing discriminant info for member type %s", tn, m))
		}
	}
}
