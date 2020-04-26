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

// Turns out basically all builders are just an embed of the corresponding assembler.
func emitEmitNodeBuilderType_typical(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Builder struct {
			_{{ .Type | TypeSymbol }}__Assembler
		}
	`, w, adjCfg, data)
}

// Builder build and reset methods are common even when some parts of the assembler vary.
// We count on the zero value of any addntl non-common fields of the assembler being correct.
func emitNodeBuilderMethods_typical(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	doTemplate(`
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Build() ipld.Node {
			if *nb.m != schema.Maybe_Value {
				panic("invalid state: cannot call Build on an assembler that's not finished")
			}
			return nb.w
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Reset() {
			var w _{{ .Type | TypeSymbol }}
			var m schema.Maybe
			*nb = _{{ .Type | TypeSymbol }}__Builder{_{{ .Type | TypeSymbol }}__Assembler{w: &w, m: &m}}
		}
	`, w, adjCfg, data)
}

// emitNodeAssemblerType_scalar emits a NodeAssembler that's typical for a scalar.
// Types that are recursive tend to have more state and custom stuff, so won't use this
// (although the 'm' and 'w' variable names may still be presumed universally).
func emitNodeAssemblerType_scalar(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe
		}

		func (na *_{{ .Type | TypeSymbol }}__Assembler) reset() {}
	`, w, adjCfg, data)
}

/*
	Some things that might make it easier to DRY stuff, if they were standard in the 'data' obj:

		- KindUpper
			- e.g. makes it possible to hoist out 'AssignNull', which needs to refer to a kind-particular mixin type.
			- also usable for to make 'AssignNode' hoistable for scalars.
				- works purely textually for scalars, conveniently.
				- maps and lists would need to branch entirely for the bottom half of the method.
		- IsRepr
			- ...?  Somewhat unsure on this one; many different ways to cut this.
			- Would be used as `{{if .IsRepr}}Repr{{end}}` in some cases, and `{{if .IsRepr}}__Repr{{end}}` in others...
			  which is viable, but somewhat disconcerting?  I dunno, maybe it's fine.
			- Also would be sometimes used as `{{if .IsRepr}}.Repr{{end}}`, in the middle of some help and error texts.

*/
