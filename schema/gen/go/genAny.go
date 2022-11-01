package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

type anyGenerator struct {
	AdjCfg  *AdjunctCfg
	PkgName string
	Type    *schema.TypeAny
}

func (anyGenerator) IsRepr() any { return false } // hint used in some generalized templates.

// --- native content and specializations --->

func (g anyGenerator) EmitNativeType(w io.Writer) {
	doTemplate(`
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}
		type _{{ .Type | TypeSymbol }} struct{ x datamodel.Node }
	`, w, g.AdjCfg, g)
}
func (g anyGenerator) EmitNativeAccessors(w io.Writer) {
}
func (g anyGenerator) EmitNativeBuilder(w io.Writer) {
}
func (g anyGenerator) EmitNativeMaybe(w io.Writer) {
	emitNativeMaybe(w, g.AdjCfg, g)
}

// --- type info --->

func (g anyGenerator) EmitTypeConst(w io.Writer) {
}

// --- TypedNode interface satisfaction --->

func (g anyGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g anyGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	emitTypicalTypedNodeMethodRepresentation(w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g anyGenerator) EmitNodeType(w io.Writer) {
}
func (g anyGenerator) EmitNodeTypeAssertions(w io.Writer) {
}
func (g anyGenerator) EmitNodeMethodPrototype(w io.Writer) {
	emitNodeMethodPrototype_typical(w, g.AdjCfg, g)
}
func (g anyGenerator) EmitNodePrototypeType(w io.Writer) {
	emitNodePrototypeType_typical(w, g.AdjCfg, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g anyGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return anyBuilderGenerator{
		g.AdjCfg,
		g.PkgName,
		g.Type,
	}
}

type anyBuilderGenerator struct {
	AdjCfg  *AdjunctCfg
	PkgName string
	Type    *schema.TypeAny
}

func (anyBuilderGenerator) IsRepr() any { return false } // hint used in some generalized templates.

func (anyBuilderGenerator) Kind() datamodel.Kind { return datamodel.Kind_Invalid }

func (g anyBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	emitEmitNodeBuilderType_typical(w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	emitNodeBuilderMethods_typical(w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	emitNodeAssemblerType_scalar(w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignNull() error {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignBool(bool) error {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignNode(v datamodel.Node) error {
			na.w.x = v
			*na.m = schema.Maybe_Value
			return nil
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// Nothing needed here for any kinds.
}

func (g anyBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) BeginMap(sizeHint int64) (datamodel.MapAssembler, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) BeginList(sizeHint int64) (datamodel.ListAssembler, error) {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignInt(int64) error {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignFloat(float64) error {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignString(string) error {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignBytes([]byte) error {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) AssignLink(datamodel.Link) error {
			panic("not implemented")
		}
	`, w, g.AdjCfg, g)
}
func (g anyBuilderGenerator) EmitNodeAssemblerMethodPrototype(w io.Writer) {
	doTemplate(`
		func (_{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Assembler) Prototype() datamodel.NodePrototype {
			return _{{ .Type | TypeSymbol }}__{{ if .IsRepr }}Repr{{end}}Prototype{}
		}
	`, w, g.AdjCfg, g)
}
