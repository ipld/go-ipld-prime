package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type linkGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.LinkTraits
	PkgName string
	Type    schema.TypeLink
}

// --- native content and specializations --->

func (g linkGenerator) EmitNativeType(w io.Writer) {
	emitNativeType_scalar(w, g.AdjCfg, g)
}
func (g linkGenerator) EmitNativeAccessors(w io.Writer) {
	emitNativeAccessors_scalar(w, g.AdjCfg, g)
}
func (g linkGenerator) EmitNativeBuilder(w io.Writer) {
	emitNativeBuilder_scalar(w, g.AdjCfg, g)
}

func (g linkGenerator) EmitNativeMaybe(w io.Writer) {
	emitNativeMaybe(w, g.AdjCfg, g)
}

// --- type info --->

func (g linkGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode linkerface satisfaction --->

func (g linkGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g linkGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	emitTypicalTypedNodeMethodRepresentation(w, g.AdjCfg, g)
}

// --- Node linkerface satisfaction --->

func (g linkGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}
func (g linkGenerator) EmitNodeTypeAssertions(w io.Writer) {
	emitNodeTypeAssertions_typical(w, g.AdjCfg, g)
}
func (g linkGenerator) EmitNodeMethodAsLink(w io.Writer) {
	emitNodeMethodAsKind_scalar(w, g.AdjCfg, g)
}
func (g linkGenerator) EmitNodeMethodStyle(w io.Writer) {
	emitNodeMethodStyle_typical(w, g.AdjCfg, g)
}
func (g linkGenerator) EmitNodeStyleType(w io.Writer) {
	emitNodeStyleType_typical(w, g.AdjCfg, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g linkGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return linkBuilderGenerator{
		g.AdjCfg,
		mixins.LinkAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__",
		},
		g.PkgName,
		g.Type,
	}
}

type linkBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.LinkAssemblerTraits
	PkgName string
	Type    schema.TypeLink
}

func (linkBuilderGenerator) IsRepr() bool { return false } // hint used in some generalized templates.

func (g linkBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	emitEmitNodeBuilderType_typical(w, g.AdjCfg, g)
}
func (g linkBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	emitNodeBuilderMethods_typical(w, g.AdjCfg, g)
}
func (g linkBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	emitNodeAssemblerType_scalar(w, g.AdjCfg, g)
}
func (g linkBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	emitNodeAssemblerMethodAssignNull_scalar(w, g.AdjCfg, g)
}
func (g linkBuilderGenerator) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	emitNodeAssemblerMethodAssignKind_scalar(w, g.AdjCfg, g)
}
func (g linkBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	emitNodeAssemblerMethodAssignNode_scalar(w, g.AdjCfg, g)
}
func (g linkBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// Nothing needed here for link kinds.
}
