package typed

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/schema"
)

// ErrNoSuchField may be returned from traversal functions on the Node
// interface when a field is requested which doesn't exist; it may also be
// returned during MapBuilder
type ErrNoSuchField struct {
	Type schema.Type

	FieldName string
}

func (e ErrNoSuchField) Error() string {
	return fmt.Sprintf("no such field: %s.%s", e.Type.Name(), e.FieldName)
}
