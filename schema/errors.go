package schema

import (
	"fmt"
)

// ErrNoSuchField may be returned from lookup functions on the Node
// interface when a field is requested which doesn't exist, or from Insert
// on a MapBuilder when a key doesn't match a field name in the structure.
type ErrNoSuchField struct {
	Type Type

	FieldName string
}

func (e ErrNoSuchField) Error() string {
	if e.Type == nil {
		return fmt.Sprintf("no such field: {typeinfomissing}.%s", e.FieldName)
	}
	return fmt.Sprintf("no such field: %s.%s", e.Type.Name(), e.FieldName)
}
