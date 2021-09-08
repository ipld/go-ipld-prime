package schemadmt

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/schema"
)

// Compile transforms a schema in DMT form into a TypeSystem.
//
// Note that this API is EXPERIMENTAL and will likely change.
// It is also unfinished and buggy.
func Compile(ts *schema.TypeSystem, node *Schema) error {
	for _, name := range node.Types.Keys {
		defn := node.Types.Values[name]

		name := schema.TypeName(name)
		switch {
		// Scalar types without parameters.
		case defn.TypeBool != nil:
			ts.Accumulate(schema.SpawnBool(name))
		case defn.TypeString != nil:
			ts.Accumulate(schema.SpawnString(name))
		case defn.TypeBytes != nil:
			ts.Accumulate(schema.SpawnBytes(name))
		case defn.TypeInt != nil:
			ts.Accumulate(schema.SpawnInt(name))
		case defn.TypeFloat != nil:
			ts.Accumulate(schema.SpawnFloat(name))

		case defn.TypeMap != nil:
			typ := defn.TypeMap
			if typ.ValueType.InlineDefn != nil {
				return fmt.Errorf("TODO: support anonymous types in schema package")
			}
			switch {
			case typ.Representation.MapRepresentation_Map != nil:
				// default behavior
			default:
				return fmt.Errorf("TODO: support non-default map repr in schema package")
			}
			ts.Accumulate(schema.SpawnMap(name,
				schema.TypeName(typ.KeyType),
				schema.TypeName(*typ.ValueType.TypeName),
				typ.ValueNullable,
			))
		case defn.TypeStruct != nil:
			typ := defn.TypeStruct
			var fields []schema.StructField
			for _, fname := range typ.Fields.Keys {
				field := typ.Fields.Values[fname]
				if field.Type.InlineDefn != nil {
					return fmt.Errorf("TODO: support anonymous types in schema package")
				}
				fields = append(fields, schema.SpawnStructField(fname,
					schema.TypeName(*field.Type.TypeName),
					field.Optional, field.Nullable,
				))
			}
			ts.Accumulate(schema.SpawnStruct(name,
				fields,
				nil, // TODO
			))
		// TODO: unions etc
		default:
		}
	}

	if errs := ts.ValidateGraph(); errs != nil {
		for _, err := range errs {
			return err
		}
	}
	return nil
}
