package selector

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
)

// SelectorSpec is a specification for a selector that can build
// a selector ipld.Node or an actual parsed Selector
type SelectorSpec interface {
	Node() ipld.Node
	Selector() (Selector, error)
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
	ExploreRecursive(maxDepth int, sequence SelectorSpec) SelectorSpec
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
	Delete(k string)
}

type selectorSpecBuilder struct {
	fnb fluent.NodeBuilder
}

type selectorSpec struct {
	n ipld.Node
}

func (ss selectorSpec) Node() ipld.Node {
	return ss.n
}

func (ss selectorSpec) Selector() (Selector, error) {
	return ParseSelector(ss.n)
}

// NewSelectorSpecBuilder creates a SelectorSpecBuilder from an underlying ipld NodeBuilder
func NewSelectorSpecBuilder(nb ipld.NodeBuilder) SelectorSpecBuilder {
	fnb := fluent.WrapNodeBuilder(nb)
	return &selectorSpecBuilder{fnb}
}

func (ssb *selectorSpecBuilder) ExploreRecursiveEdge() SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreRecursiveEdgeKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreRecursive(maxDepth int, sequence SelectorSpec) SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreRecursiveKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(maxDepthKey), vnb.CreateInt(maxDepth))
				mb.Insert(knb.CreateString(sequenceKey), sequence.Node())
			}))
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreAll(next SelectorSpec) SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreAllKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(nextSelectorKey), next.Node())
			}))
		}),
	}
}
func (ssb *selectorSpecBuilder) ExploreIndex(index int, next SelectorSpec) SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreIndexKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(indexKey), vnb.CreateInt(index))
				mb.Insert(knb.CreateString(nextSelectorKey), next.Node())
			}))
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreRange(start int, end int, next SelectorSpec) SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreRangeKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(startKey), vnb.CreateInt(start))
				mb.Insert(knb.CreateString(endKey), vnb.CreateInt(end))
				mb.Insert(knb.CreateString(nextSelectorKey), next.Node())
			}))
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreUnion(members ...SelectorSpec) SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreUnionKey), vnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
				for _, member := range members {
					lb.Append(member.Node())
				}

			}))
		}),
	}
}

func (ssb *selectorSpecBuilder) ExploreFields(specBuilder ExploreFieldsSpecBuildingClosure) SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreFieldsKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(fieldsKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					specBuilder(exploreFieldsSpecBuilder{mb, knb})
				}))
			}))
		}),
	}
}

func (ssb *selectorSpecBuilder) Matcher() SelectorSpec {
	return selectorSpec{
		ssb.fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
		}),
	}
}

type exploreFieldsSpecBuilder struct {
	mb  fluent.MapBuilder
	knb fluent.NodeBuilder
}

func (efsb exploreFieldsSpecBuilder) Insert(field string, s SelectorSpec) {
	efsb.mb.Insert(efsb.knb.CreateString(field), s.Node())
}

func (efsb exploreFieldsSpecBuilder) Delete(field string) {
	efsb.mb.Delete(efsb.knb.CreateString(field))
}
