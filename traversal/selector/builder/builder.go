package builder

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	selector "github.com/ipld/go-ipld-prime/traversal/selector"
)

// SelectorSpec is a specification for a selector that can build
// a selector ipld.Node or an actual parsed Selector
type SelectorSpec interface {
	Node() ipld.Node
	Selector() (selector.Selector, error)
}

// SelectorSpecBuilder is a utility interface to build selector ipld nodes
// quickly.
//
// It serves two purposes:
// 1. Save the user of go-ipld-prime time and mental overhead with an easy
// interface for making selector nodes in much less code without having to remember
// the selector sigils
// 2. Provide a level of protection from selector schema changes, at least in terms
// of naming, if not structure
type SelectorSpecBuilder interface {
	ExploreRecursiveEdge() SelectorSpec
	ExploreRecursive(limit selector.RecursionLimit, sequence SelectorSpec) SelectorSpec
	ExploreUnion(...SelectorSpec) SelectorSpec
	ExploreAll(next SelectorSpec) SelectorSpec
	ExploreIndex(index int, next SelectorSpec) SelectorSpec
	ExploreRange(start int, end int, next SelectorSpec) SelectorSpec
	ExploreFields(ExploreFieldsSpecBuildingClosure) SelectorSpec
	Matcher() SelectorSpec
}

// ExploreFieldsSpecBuildingClosure is a function that provided to SelectorSpecBuilder's
// ExploreFields method that assembles the fields map in the selector using
// an ExploreFieldsSpecBuilder
type ExploreFieldsSpecBuildingClosure func(ExploreFieldsSpecBuilder)

// ExploreFieldsSpecBuilder is an interface for assemble the map of fields to
// selectors in ExploreFields
type ExploreFieldsSpecBuilder interface {
	Insert(k string, v SelectorSpec)
}

type selectorSpecBuilder struct {
	ns ipld.NodeStyle
}

type selectorSpec struct {
	n ipld.Node
}

func (ss selectorSpec) Node() ipld.Node {
	return ss.n
}

func (ss selectorSpec) Selector() (selector.Selector, error) {
	return selector.ParseSelector(ss.n)
}

// NewSelectorSpecBuilder creates a SelectorSpecBuilder which will store
// data in the format determined by the given ipld.NodeStyle.
func NewSelectorSpecBuilder(ns ipld.NodeStyle) SelectorSpecBuilder {
	return &selectorSpecBuilder{ns}
}

func (ssb *selectorSpecBuilder) ExploreRecursiveEdge() SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapNodeAssembler) {})
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreRecursive(limit selector.RecursionLimit, sequence SelectorSpec) SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreRecursive).CreateMap(2, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Limit).CreateMap(1, func(na fluent.MapNodeAssembler) {
					switch limit.Mode() {
					case selector.RecursionLimit_Depth:
						na.AssembleDirectly(selector.SelectorKey_LimitDepth).AssignInt(limit.Depth())
					case selector.RecursionLimit_None:
						na.AssembleDirectly(selector.SelectorKey_LimitNone).CreateMap(0, func(na fluent.MapNodeAssembler) {})
					default:
						panic("Unsupported recursion limit type")
					}
				})
				na.AssembleDirectly(selector.SelectorKey_Sequence).AssignNode(sequence.Node())
			})
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreAll(next SelectorSpec) SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Next).AssignNode(next.Node())
			})
		}),
	}
}
func (ssb *selectorSpecBuilder) ExploreIndex(index int, next SelectorSpec) SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreIndex).CreateMap(2, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Index).AssignInt(index)
				na.AssembleDirectly(selector.SelectorKey_Next).AssignNode(next.Node())
			})
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreRange(start int, end int, next SelectorSpec) SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreRange).CreateMap(3, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Start).AssignInt(start)
				na.AssembleDirectly(selector.SelectorKey_End).AssignInt(end)
				na.AssembleDirectly(selector.SelectorKey_Next).AssignNode(next.Node())
			})
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreUnion(members ...SelectorSpec) SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreUnion).CreateList(len(members), func(na fluent.ListNodeAssembler) {
				for _, member := range members {
					na.AssembleValue().AssignNode(member.Node())
				}
			})
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreFields(specBuilder ExploreFieldsSpecBuildingClosure) SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreFields).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Fields).CreateMap(-1, func(na fluent.MapNodeAssembler) {
					specBuilder(exploreFieldsSpecBuilder{na})
				})
			})
		}),
	}
}

func (ssb *selectorSpecBuilder) Matcher() SelectorSpec {
	return selectorSpec{
		fluent.MustBuildMap(ssb.ns, 1, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapNodeAssembler) {})
		}),
	}
}

type exploreFieldsSpecBuilder struct {
	na fluent.MapNodeAssembler
}

func (efsb exploreFieldsSpecBuilder) Insert(field string, s SelectorSpec) {
	efsb.na.AssembleDirectly(field).AssignNode(s.Node())
}
