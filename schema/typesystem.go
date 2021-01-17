package schema

type TypeSystem struct {
	// Mind the key type here: TypeReference, not TypeName.
	// The key might be a computed anon "name" which is not actually a valid type name itself.
	types map[TypeReference]Type

	// List of types, retained in the original order they were specified,
	// including only those which are named (not any computed anonymous types).
	// This is kept so we can do any listing in the order the user expects,
	// report any errors during rule validation in the same order as the input, etc.
	list []Type

	// List of anonymous types, in sorted order simply going by the reference as a string.
	anonTypes []Type
}

// AllTypes returns a slice of every Type in this TypeSystem.
// This slice includes both named types and anonymous types.
// For named types, the order in which they were specified
// then the TypeSystem was compiled is preserved.
func (ts *TypeSystem) AllTypes() []Type {
	v := make([]Type, len(ts.types))
	copy(ts.list, v)
	copy(ts.anonTypes, v[len(ts.list):])
	return v
}

func (ts *TypeSystem) GetType(name TypeReference) Type {
	return ts.types[name]
}
