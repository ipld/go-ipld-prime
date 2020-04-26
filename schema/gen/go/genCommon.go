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

func emitNodeAssemblerMethodAssignNull_scalar(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignNull() error {
			switch *na.m {
			case allowNull:
				*na.m = schema.Maybe_Null
				return nil
			case schema.Maybe_Absent:
				return mixins.{{ .ReprKind.String | title }}Assembler{"{{ .PkgName }}.{{ .TypeName }}{{ if .IsRepr }}.Repr{{end}}"}.AssignNull()
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			}
			panic("unreachable")
		}
	`, w, adjCfg, data)
}

// almost the same as the variant for scalars, but also has to check for midvalue state.
func emitNodeAssemblerMethodAssignNull_recursive(w io.Writer, adjCfg *AdjunctCfg, data interface{}) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignNull() error {
			switch *na.m {
			case allowNull:
				*na.m = schema.Maybe_Null
				return nil
			case schema.Maybe_Absent:
				return mixins.{{ .ReprKind.String | title }}Assembler{"{{ .PkgName }}.{{ .TypeName }}{{ if .IsRepr }}.Repr{{end}}"}.AssignNull()
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			case midvalue:
				panic("invalid state: cannot assign null into an assembler that's already begun working on recursive structures!")
			}
			panic("unreachable")
		}
	`, w, adjCfg, data)
}
