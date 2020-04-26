package gengo

import (
	"io"
)

/*
	This file is full of "typical" templates.
	They may not be used by *every* type and representation,
	but if they're extracted here, they're at least used by *many*.
*/

// emitNativeMaybe turns out to be completely agnostic to pretty much everything;
// it doesn't vary by kind at all, and has never yet ended up needing specialization.
func emitNativeMaybe(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Maybe struct {
			m schema.Maybe
			v {{if not (MaybeUsesPtr .Type) }}_{{end}}{{ .Type | TypeSymbol }}
		}
		type Maybe{{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}__Maybe

		func (m Maybe{{ .Type | TypeSymbol }}) IsNull() bool {
			return m.m == schema.Maybe_Null
		}
		func (m Maybe{{ .Type | TypeSymbol }}) IsUndefined() bool {
			return m.m == schema.Maybe_Absent
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Exists() bool {
			return m.m == schema.Maybe_Value
		}
		func (m Maybe{{ .Type | TypeSymbol }}) AsNode() ipld.Node {
			switch m.m {
				case schema.Maybe_Absent:
					return ipld.Undef
				case schema.Maybe_Null:
					return ipld.Null
				case schema.Maybe_Value:
					return {{if not (MaybeUsesPtr .Type) }}&{{end}}m.v
				default:
					panic("unreachable")
			}
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Must() {{ .Type | TypeSymbol }} {
			if !m.Exists() {
				panic("unbox of a maybe rejected")
			}
			return {{if not (MaybeUsesPtr .Type) }}&{{end}}m.v
		}
	`, w, adjCfg, data)
}

// emitTypicalTypedNodeMethodRepresentation does... what it says on the tin.
//
// For most types, the way to get the representation node pointer doesn't
// textually depend on either the node implementation details nor what the representation strategy is,
// or really much at all for that matter.
// It only depends on that they have the same structure, so this cast works.
//
// Most (all?) types can use this.  However, it's here rather in the mixins, for two reasons:
// one, it still seems possible to imagine we'll have a type someday for which this pattern won't hold;
// and two, mixins are also used in the repr generators, and it wouldn't be all sane for this method to end up also on reprs.
func emitTypicalTypedNodeMethodRepresentation(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Representation() ipld.Node {
			return (*_{{ .Type | TypeSymbol }}__Repr)(n)
		}
	`, w, adjCfg, data)
}
