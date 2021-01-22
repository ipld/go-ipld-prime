/*
	The schema/compiler package contains concrete implementations of the
	interfaces in the schema package which are used to describe IPLD Schemas,
	and it also provides a Compiler type which is used to construct them.
*/
package schema

import (
	"fmt"
)

// Compiler creates new TypeSystem instances.
// Methods are called on a Compiler instance to add types to the set,
// and when done, the Compile method is called, which can return
// either a list of error values or a new TypeSystem.
//
// Users don't usually use Compiler themselves,
// and this API isn't meant to be especially user-friendly.
// It's better to write IPLD Schemas using the DSL,
// parse and transpile that into the standard DMT format,
// and then read that with `schema/dmt` package and use the `dmt.Compile` feature.
// This lets you spend more time with the human-readable syntax and DMT format,
// which in addition to being better suited for documentation and review,
// is also usable with other IPLD tools and IPLD implementations in other languages.
// (Inside, the `dmt.Compile` feature uses this Compiler for you.)
//
// On error handling:
// Since several sorts of error can't be checked until the whole group of types has been stated
// (for example, referential completeness checks),
// almost none of the methods on Compiler return errors as they go.
// All errors will just be reported altogether at once at the end, when Compile is called.
// Some extremely obvious errors, like trying to use the same TypeName twice, will cause a panic immediately.
// The rule for errors that are raised as panics is that they must have been already avoided if the data were coming from the schemadmt package.
// (E.g., if something could be invalidly sent to the Compiler twice, but was a map key in the schemadmt and so already checked as unique, that's panic-worthy here.
// But if repeats of some identifier are invalid but would a list when expressed in the schemadmt, that's *not* allowed to panic here.)
//
// On immutability:
// The TypeSystem returned by a successful Compile call will be immutable.
// Many methods on the Compiler type are structured to accept data in a way that works towards this immutability.
// In particular, many methods on Compiler take arguments which are "carrier types" for segments of immutable data,
// and these "carrier types" are produced by constructor functions.
// For one example of this pattern, see the interplay of compiler.TypeStruct() and MakeStructFieldList().
//
// On code organization:
// Several methods are attached to the Compiler type but don't actually take it as a parameter.
// (All these methods have the name prefix "Make*".)
// These methods are constructors for various intermediate values needed to feed information into the compiler.
// These are attached to the Compiler type purely for organization of the godoc,
// so that they don't clutter up the package with functions that users should never be expected to use.
type Compiler struct {
	// ... and if you're wondering why this type is exported at all?
	//  Well, arguably, it's useful to be able to construct these values without going through the dmt.
	//  At the end of the day, though?  Honestly, import cycle breaking.  This was not the first choice.
	// An implementation which wraps the schemadmt package to make it fit the schema interfaces was the first choice
	//  because it would've saved a *lot* of work (it would've removed the need for this compiler system entirely, among other things);
	//  but that doesn't fly, because the generated nodes have to implement schema.TypedNode, and that interface refers to other types like schema.Type
	//   (and even if all of those are interfaces, you still can't have conformational-equality span multiple interface types in golang).
	//  And that would make an import cycle if we tried to put types wrapping the dmt types into the schema package.  Whoops.
	// A separation of this feature into its own "compiler" package was also attempted; this too did not work out well.
	//  (The main reason to desire this would be clarity and grouping.  The reading interfaces and the creation stuff are for different audiences.)
	//  We want immutability in almost all of these values, and in golang, the only way to approach immutability is with unexported symbols and package boundaries;
	//  there aren't many practical approaches which would allow separating the reading and the creating parts of handling the same structure into two different packages.
	//  One approach that does allow such a split is to simply define all reading via interfaces.  This works and is not wildly unusual in golang.
	//   However, to apply this in practice in this case would result in many interfaces which are just shells around exactly one implementation;
	//   and it would also suggest various implementations are expected, which... is simply not the case.
	//   This would still overall be viable, but the result would be pushing a bit into "strange".
	// So, here we are.
	//  Compiler is a type in the schema package.  And we've attempted to attach *all* operations relating to creating schema data to it for purposes of grouping and clarity.

	// ts gathers all the in-progress types (including anonymous ones),
	// and is eventually the value we return (if Compile is ultimately successful).
	// We insert into this blindly as we go, and check everything for consistency at the end;
	// if those logical checks flunk, we don't allow any reference to it to escape.
	// This is nil'd after any Compile, so when we give a reference to it away,
	// it's immutable from there on out.
	ts *TypeSystem
}

func (c *Compiler) Init() {
	c.ts = &TypeSystem{
		map[TypeReference]Type{},
		nil,
		nil,
	}
}

func (c *Compiler) Compile() (*TypeSystem, error) {
	panic("TODO")
}

func (c *Compiler) MustCompile() *TypeSystem {
	ts, err := c.Compile()
	if err != nil {
		panic(err)
	}
	return ts
}

func (c *Compiler) addType(t Type) {
	c.mustHaveNameFree(t.Name())
	c.ts.types[TypeReference(t.Name())] = t
	c.ts.list = append(c.ts.list, t)
}
func (c *Compiler) addAnonType(t Type) {
	c.ts.types[TypeReference(t.Name())] = t // FIXME it's... probably a bug that the Type.Name() method doesn't return a TypeReference.  Yeah, it definitely is.  TypeMap and TypeList should have their own name field internally be TypeReference, too, because it's true.  wonder if we should have separate methods on the Type interface for this.  would probably be a usability trap to do so, though (too many user printfs would use the Name function and get blanks and be surprised).
}

func (c *Compiler) mustHaveNameFree(name TypeName) {
	if _, exists := c.ts.types[TypeReference(name)]; exists {
		panic(fmt.Errorf("type name %q already used", name))
	}
}

//go:generate sed -i /---/q compiler_carriers.go

func (c *Compiler) TypeBool(name TypeName) {
	c.addType(&TypeBool{c.ts, name})
}

func (c *Compiler) TypeString(name TypeName) {
	c.addType(&TypeString{c.ts, name})
}

func (c *Compiler) TypeBytes(name TypeName) {
	c.addType(&TypeBytes{c.ts, name})
}

func (c *Compiler) TypeInt(name TypeName) {
	c.addType(&TypeInt{c.ts, name})
}

func (c *Compiler) TypeFloat(name TypeName) {
	c.addType(&TypeFloat{c.ts, name})
}

func (c *Compiler) TypeLink(name TypeName, expectedTypeRef TypeName) {
	c.addType(&TypeLink{c.ts, name, expectedTypeRef})
}

func (c *Compiler) TypeStruct(name TypeName, fields structFieldList, rstrat StructRepresentation) {
	t := TypeStruct{
		ts:        c.ts,
		name:      name,
		fields:    fields.x, // it's safe to take this directly because the carrier type means a reference to this slice has never been exported.
		fieldsMap: make(map[StructFieldName]*StructField, len(fields.x)),
		rstrat:    rstrat,
	}
	c.addType(&t)
	for i, f := range fields.x {
		// duplicate names are rejected with a *panic* here because we expect these to already be unique (if this data is coming from the dmt, these were map keys there).
		if _, exists := t.fieldsMap[f.name]; exists {
			panic(fmt.Errorf("type %q already has field named %q", t.name, f.name))
		}
		t.fieldsMap[f.name] = &fields.x[i]
		fields.x[i].parent = &t
	}
}

//go:generate quickimmut -output=compiler_carriers.go -attach=Compiler list StructField

func (Compiler) MakeStructField(name StructFieldName, typ TypeReference, optional, nullable bool) StructField {
	return StructField{nil, name, typ, optional, nullable}
}

func (Compiler) MakeStructRepresentation_Map(fieldDetails structFieldNameStructRepresentation_Map_FieldDetailsMap) StructRepresentation {
	return StructRepresentation_Map{nil, fieldDetails.x}
}

//go:generate quickimmut -output=compiler_carriers.go -attach=Compiler map StructFieldName StructRepresentation_Map_FieldDetails

func (c *Compiler) TypeMap(name TypeName, keyTypeRef TypeName, valueTypeRef TypeReference, valueNullable bool, rstrat MapRepresentation) {
	c.addType(&TypeMap{c.ts, name, keyTypeRef, valueTypeRef, valueNullable, rstrat})
}

func (Compiler) MakeMapRepresentation_Stringpairs(innerDelim string, entryDelim string) MapRepresentation {
	return MapRepresentation_Stringpairs{innerDelim, entryDelim}
}

func (c *Compiler) TypeList(name TypeName, valueTypeRef TypeReference, valueNullable bool) {
	c.addType(&TypeList{c.ts, name, valueTypeRef, valueNullable})
}

func (c *Compiler) TypeUnion(name TypeName, members typeNameList, rstrat UnionRepresentation) {
	t := TypeUnion{
		ts:      c.ts,
		name:    name,
		members: members.x, // it's safe to take this directly because the carrier type means a reference to this slice has never been exported.
		rstrat:  rstrat,
	}
	c.addType(&t)
	// note! duplicate member names *not* rejected at this moment -- that's a job for the validation phase.
	//  this is an interesting contrast to how when building structs, dupe field names may be rejected proactively:
	//   the difference is, member names were a list in the dmt form too, so it's important we format a nice error rather than panic if there was invalid data there.
}

//go:generate quickimmut -output=compiler_carriers.go -attach=Compiler list TypeName

func (Compiler) MakeUnionRepresentation_Keyed(discriminantTable stringTypeNameMap) UnionRepresentation {
	return &UnionRepresentation_Keyed{nil, discriminantTable.x}
}

//go:generate quickimmut -output=compiler_carriers.go -attach=Compiler map string TypeName
