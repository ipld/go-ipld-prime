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

// FUTURE/DRY: unclear how this should be extracted.  (if?  yes.  how?  dunno.)
// One on hand, it seems tempting to put it in the mixins.
// On the other hand, it's not at all sane for this to end up also embedded in reprs, so... no.
// Bifrucate the mixins?  Maybe, but is that going to create more work overall or less?
func emitTypicalTypedNodeMethodRepresentation(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	// Perhaps surprisingly, the way to get the representation node pointer
	//  tends not to actually depend on either the node implementation details nor what the representation strategy is.
	//   It only depends on that they have the same structure -- and there's no textual appearance of the details here.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Representation() ipld.Node {
			return (*_{{ .Type | TypeSymbol }}__Repr)(n)
		}
	`, w, adjCfg, data)
}
