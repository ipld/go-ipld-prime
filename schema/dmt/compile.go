package schemadmt

import (
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/compiler"
)

// This code is broken up into a bunch of individual 'compile' methods,
//  each attached to the type that's their input information.
// However, many of them return distinct concrete types,
//  and so we've just chained it all together with switch statements;
//   creating a separate interface per result type seems just not super relevant.

func (schdmt Schema) Compile() (schema.TypeSystem, error) {
	c := &compiler.Compiler{}
	typesdmt := schdmt.FieldTypes()
	for itr := typesdmt.Iterator(); !itr.Done(); {
		tn, t := itr.Next()
		switch t2 := t.AsInterface().(type) {
		case TypeBool:
			c.TypeBool(schema.TypeName(tn.String()))
		case TypeString:
			c.TypeString(schema.TypeName(tn.String()))
		case TypeBytes:
			c.TypeBytes(schema.TypeName(tn.String()))
		case TypeInt:
			c.TypeInt(schema.TypeName(tn.String()))
		case TypeFloat:
			c.TypeFloat(schema.TypeName(tn.String()))
		case TypeLink:
			if t2.FieldExpectedType().Exists() {
				c.TypeLink(schema.TypeName(tn.String()), schema.TypeName(t2.FieldExpectedType().Must().String()))
			} else {
				c.TypeLink(schema.TypeName(tn.String()), "")
			}
		case TypeMap:
			c.TypeMap(
				schema.TypeName(tn.String()),
				schema.TypeName(t2.FieldKeyType().String()),
				t2.FieldValueType().TypeReference(),
				t2.FieldValueNullable().Bool(),
			)
			// If the field typeReference is TypeDefnInline, that needs a chance to take additional action.
			t2.FieldValueType().compile(c)
		case TypeList:
			c.TypeList(
				schema.TypeName(tn.String()),
				t2.FieldValueType().TypeReference(),
				t2.FieldValueNullable().Bool(),
			)
			// If the field typeReference is TypeDefnInline, that needs a chance to take additional action.
			t2.FieldValueType().compile(c)
		case TypeStruct:
			// Flip fields info from DMT to compiler argument format.
			fields := make([]compiler.StructField, t2.FieldFields().Length())
			for itr := t2.FieldFields().Iterator(); !itr.Done(); {
				fname, fdmt := itr.Next()
				fields = append(fields, compiler.MakeStructField(
					compiler.StructFieldName(fname.String()),
					fdmt.FieldType().TypeReference(),
					fdmt.FieldOptional().Bool(),
					fdmt.FieldNullable().Bool(),
				))
				// If the field typeReference is TypeDefnInline, that needs a chance to take additional action.
				fdmt.FieldType().compile(c)
			}
			// Flip the representaton strategy DMT to compiler argument format.
			rstrat := func() compiler.StructRepresentation {
				switch r := t2.FieldRepresentation().AsInterface().(type) {
				case StructRepresentation_Map:
					return r.compile()
				case StructRepresentation_Tuple:
					return r.compile()
				case StructRepresentation_Stringpairs:
					return r.compile()
				case StructRepresentation_Stringjoin:
					return r.compile()
				case StructRepresentation_Listpairs:
					return r.compile()
				default:
					panic("unreachable")
				}
			}()
			// Feed it all into the compiler.
			c.TypeStruct(
				schema.TypeName(tn.String()),
				compiler.MakeStructFieldList(fields...),
				rstrat,
			)
		case TypeUnion:
			// Flip members info from DMT to compiler argument format.
			members := make([]schema.TypeName, t2.FieldMembers().Length())
			for itr := t2.FieldMembers().Iterator(); !itr.Done(); {
				_, memberName := itr.Next()
				members = append(members, schema.TypeName(memberName.String()))
				// n.b. no need to check for TypeDefnInline here, because schemas don't allow those in union defns.
			}
			// Flip the representaton strategy DMT to compiler argument format.
			rstrat := func() compiler.UnionRepresentation {
				switch r := t2.FieldRepresentation().AsInterface().(type) {
				case UnionRepresentation_Keyed:
					return r.compile()
				case UnionRepresentation_Kinded:
					return r.compile()
				case UnionRepresentation_Envelope:
					return r.compile()
				case UnionRepresentation_Inline:
					return r.compile()
				case UnionRepresentation_StringPrefix:
					return r.compile()
				case UnionRepresentation_BytePrefix:
					return r.compile()
				default:
					panic("unreachable")
				}
			}()
			// Feed it all into the compiler.
			c.TypeUnion(
				schema.TypeName(tn.String()),
				compiler.MakeUnionMemberList(members...),
				rstrat,
			)
		case TypeEnum:
			panic("TODO")
		case TypeCopy:
			panic("no support for 'copy' types.  I might want to reneg on whether these are even part of the schema dmt.")
		default:
			panic("unreachable")
		}
	}
	return c.Compile()
}

// If the typeReference is TypeDefnInline, create the anonymous type and feed it to the compiler.
// It's fine if anonymous type has been seen before; we let dedup of that be handled by the compiler.
func (dmt TypeNameOrInlineDefn) compile(c *compiler.Compiler) {
	switch dmt.AsInterface().(type) {
	case TypeDefnInline:
		panic("nyi") // TODO this needs to engage in anonymous type spawning.
	}
}

func (dmt StructRepresentation_Map) compile() compiler.StructRepresentation {
	if !dmt.FieldFields().Exists() {
		return compiler.MakeStructRepresentation_Map()
	}
	fields := make([]compiler.StructRepresentation_Map_FieldDetailsEntry, dmt.FieldFields().Must().Length())
	for itr := dmt.FieldFields().Must().Iterator(); !itr.Done(); {
		fn, det := itr.Next()
		fields = append(fields, compiler.StructRepresentation_Map_FieldDetailsEntry{
			FieldName: compiler.StructFieldName(fn.String()),
			Details: compiler.StructRepresentation_Map_FieldDetails{
				Rename: func() string {
					if det.FieldRename().Exists() {
						return det.FieldRename().Must().String()
					}
					return ""
				}(),
				Implicit: nil, // TODO
			},
		})
	}
	return compiler.MakeStructRepresentation_Map(fields...)
}

func (dmt StructRepresentation_Tuple) compile() compiler.StructRepresentation {
	panic("TODO")
}

func (dmt StructRepresentation_Stringpairs) compile() compiler.StructRepresentation {
	panic("TODO")
}

func (dmt StructRepresentation_Stringjoin) compile() compiler.StructRepresentation {
	panic("TODO")
}

func (dmt StructRepresentation_Listpairs) compile() compiler.StructRepresentation {
	panic("TODO")
}

func (dmt UnionRepresentation_Keyed) compile() compiler.UnionRepresentation {
	ents := make([]compiler.UnionDiscriminantStringEntry, 0, dmt.Length())
	for itr := dmt.Iterator(); !itr.Done(); {
		k, v := itr.Next()
		ents = append(ents, compiler.UnionDiscriminantStringEntry{k.String(), schema.TypeName(v.String())})
	}
	return compiler.MakeUnionRepresentation_Keyed(compiler.MakeUnionDiscriminantStringTable(ents...))
}

func (dmt UnionRepresentation_Kinded) compile() compiler.UnionRepresentation {
	panic("TODO")
}

func (dmt UnionRepresentation_Envelope) compile() compiler.UnionRepresentation {
	panic("TODO")
}

func (dmt UnionRepresentation_Inline) compile() compiler.UnionRepresentation {
	panic("TODO")
}

func (dmt UnionRepresentation_StringPrefix) compile() compiler.UnionRepresentation {
	panic("TODO")
}

func (dmt UnionRepresentation_BytePrefix) compile() compiler.UnionRepresentation {
	panic("TODO")
}
