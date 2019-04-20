package typedeclaration

import "fmt"

type String string

// Get looks up a value in the map by key.
//
// TraverseField performs a similar function, but is the ipld.Node generic variant;
// Get takes the native typed key and explicitly returns the native typed value.
//
// This is a generated method.
func (m AnonMap__TypeStruct__fields) Get(k String) (*StructField, error) {
	v, ok := m.val[k]
	if !ok {
		return nil, fmt.Errorf("404")
	}
	return &v, nil
}

type AnonMap__TypeStruct__fields struct {
	val map[String]StructField
	ord []string
}
