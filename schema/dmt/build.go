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
	ts.Accumulate(schema.SpawnBool("Bool"))
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnFloat("Float"))
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnBytes("Bytes"))
	for _, name := range node.Types.Keys {
		defn := node.Types.Values[name]

		name := schema.TypeName(name)
		// TODO: once we support anon types, remove the ts argument.
		typ, err := spawnType(ts, name, defn)
		if err != nil {
			return err
		}
		ts.Accumulate(typ)
	}

	if errs := ts.ValidateGraph(); errs != nil {
		for _, err := range errs {
			return err
		}
	}
	return nil
}

func todoFromImplicitlyFalseBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func todoAnonTypeName(nameOrDefn TypeNameOrInlineDefn) string {
	if nameOrDefn.TypeName != nil {
		return *nameOrDefn.TypeName
	}
	defn := *nameOrDefn.InlineDefn
	switch {
	case defn.TypeDefnMap != nil:
		defn := defn.TypeDefnMap
		return fmt.Sprintf("Map__%s__%s", defn.KeyType, todoAnonTypeName(defn.ValueType))
	case defn.TypeDefnList != nil:
		defn := defn.TypeDefnList
		return fmt.Sprintf("List__%s", todoAnonTypeName(defn.ValueType))
	default:
		panic(fmt.Errorf("%#v", defn))
	}
}

func spawnType(ts *schema.TypeSystem, name schema.TypeName, defn TypeDefn) (schema.Type, error) {
	switch {
	// Scalar types without parameters.
	case defn.TypeDefnBool != nil:
		return schema.SpawnBool(name), nil
	case defn.TypeDefnString != nil:
		return schema.SpawnString(name), nil
	case defn.TypeDefnBytes != nil:
		return schema.SpawnBytes(name), nil
	case defn.TypeDefnInt != nil:
		return schema.SpawnInt(name), nil
	case defn.TypeDefnFloat != nil:
		return schema.SpawnFloat(name), nil

	case defn.TypeDefnList != nil:
		typ := defn.TypeDefnList
		if typ.ValueType.InlineDefn != nil {
			return nil, fmt.Errorf("TODO: support anonymous types in schema package")
		}
		switch {
		case typ.Representation == nil ||
			typ.Representation.ListRepresentation_List != nil:
			// default behavior
		default:
			return nil, fmt.Errorf("TODO: support non-default map repr in schema package")
		}
		return schema.SpawnList(name,
			schema.TypeName(*typ.ValueType.TypeName),
			todoFromImplicitlyFalseBool(typ.ValueNullable),
		), nil
	case defn.TypeDefnMap != nil:
		typ := defn.TypeDefnMap
		if typ.ValueType.InlineDefn != nil {
			return nil, fmt.Errorf("TODO: support anonymous types in schema package")
		}
		switch {
		case typ.Representation == nil ||
			typ.Representation.MapRepresentation_Map != nil:
			// default behavior
		default:
			return nil, fmt.Errorf("TODO: support non-default map repr in schema package")
		}
		return schema.SpawnMap(name,
			schema.TypeName(typ.KeyType),
			schema.TypeName(*typ.ValueType.TypeName),
			todoFromImplicitlyFalseBool(typ.ValueNullable),
		), nil
	case defn.TypeDefnStruct != nil:
		typ := defn.TypeDefnStruct
		var fields []schema.StructField
		for _, fname := range typ.Fields.Keys {
			field := typ.Fields.Values[fname]
			var typeName schema.TypeName
			if field.Type.TypeName != nil {
				typeName = schema.TypeName(*field.Type.TypeName)
			} else if tname := todoAnonTypeName(field.Type); ts.TypeByName(tname) == nil {
				typeName = schema.TypeName(tname)
				// Note that TypeDefn and InlineDefn aren't the same enum.
				anonDefn := TypeDefn{
					TypeDefnMap:  field.Type.InlineDefn.TypeDefnMap,
					TypeDefnList: field.Type.InlineDefn.TypeDefnList,
					TypeDefnLink: field.Type.InlineDefn.TypeDefnLink,
				}
				anonType, err := spawnType(ts, typeName, anonDefn)
				if err != nil {
					return nil, err
				}
				ts.Accumulate(anonType)
			} else {
				typeName = schema.TypeName(tname)
			}
			fields = append(fields, schema.SpawnStructField(fname,
				schema.TypeName(typeName),
				todoFromImplicitlyFalseBool(field.Optional),
				todoFromImplicitlyFalseBool(field.Nullable),
			))
		}
		return schema.SpawnStruct(name,
			fields,
			nil, // TODO: struct repr
		), nil
	case defn.TypeDefnUnion != nil:
		typ := defn.TypeDefnUnion
		var members []schema.TypeName
		for _, member := range typ.Members {
			if member.TypeName != nil {
				members = append(members, schema.TypeName(*member.TypeName))
			} else {
				panic("TODO: inline union members")
			}
		}
		return schema.SpawnUnion(name,
			members,
			nil, // TODO: union repr
		), nil
	case defn.TypeDefnEnum != nil:
		typ := defn.TypeDefnEnum
		return schema.SpawnEnum(name,
			typ.Members,
			nil, // TODO: enum repr
		), nil
	default:
		panic(fmt.Errorf("%#v", defn))
	}
}
