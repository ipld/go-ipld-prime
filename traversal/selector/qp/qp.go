package qp

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	selector "github.com/ipld/go-ipld-prime/traversal/selector"
)

func exploreRecursiveEdge(ma ipld.MapAssembler) {
	qp.MapEntry(ma, selector.SelectorKey_ExploreRecursiveEdge, qp.Map(0, func(ipld.MapAssembler) {}))
}

func ExploreRecusiveEdge(pt ipld.NodePrototype) (ipld.Node, error) {
	return qp.BuildMap(pt, 1, exploreRecursiveEdge)
}

func ExploreRecursiveEdge() qp.Assemble {
	return qp.Map(1, exploreRecursiveEdge)
}

func exploreRecursive(limit selector.RecursionLimit, sequence qp.Assemble) func(ipld.MapAssembler) {
	return func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreRecursive, qp.Map(2, func(ma ipld.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Limit, qp.Map(1, func(ma ipld.MapAssembler) {
				switch limit.Mode() {
				case selector.RecursionLimit_Depth:
					qp.MapEntry(ma, selector.SelectorKey_LimitDepth, qp.Int(limit.Depth()))
				case selector.RecursionLimit_None:
					qp.MapEntry(ma, selector.SelectorKey_LimitNone, qp.Map(0, func(na ipld.MapAssembler) {}))
				default:
					panic("Unsupported recursion limit type")
				}
			}))
			qp.MapEntry(ma, selector.SelectorKey_Sequence, sequence)
		}))
	}
}

func BuildExploreRecursive(np ipld.NodePrototype, limit selector.RecursionLimit, sequence qp.Assemble) (ipld.Node, error) {
	return qp.BuildMap(np, 1, exploreRecursive(limit, sequence))
}

func ExploreRecursive(limit selector.RecursionLimit, sequence qp.Assemble) qp.Assemble {
	return qp.Map(1, exploreRecursive(limit, sequence))
}

func exploreAll(next qp.Assemble) func(ipld.MapAssembler) {
	return func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreAll, qp.Map(1, func(ma ipld.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Next, next)
		}))
	}
}

func ExploreAll(next qp.Assemble) qp.Assemble {
	return qp.Map(1, exploreAll(next))
}

func BuildExploreAll(np ipld.NodePrototype, next qp.Assemble) (ipld.Node, error) {
	return qp.BuildMap(np, 1, exploreAll(next))
}

func exploreIndex(index int64, next qp.Assemble) func(ma ipld.MapAssembler) {
	return func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreIndex, qp.Map(2, func(ma ipld.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Index, qp.Int(index))
			qp.MapEntry(ma, selector.SelectorKey_Next, next)
		}))
	}
}

func BuildExploreIndex(np ipld.NodePrototype, index int64, next qp.Assemble) (ipld.Node, error) {
	return qp.BuildMap(np, 1, exploreIndex(index, next))
}

func ExploreIndex(index int64, next qp.Assemble) qp.Assemble {
	return qp.Map(1, exploreIndex(index, next))
}

func exploreRange(start, end int64, next qp.Assemble) func(ipld.MapAssembler) {
	return func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreRange, qp.Map(3, func(ma ipld.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Start, qp.Int(start))
			qp.MapEntry(ma, selector.SelectorKey_End, qp.Int(end))
			qp.MapEntry(ma, selector.SelectorKey_Next, next)
		}))
	}
}

func BuildExploreRange(np ipld.NodePrototype, start, end int64, next qp.Assemble) (ipld.Node, error) {
	return qp.BuildMap(np, 1, exploreRange(start, end, next))
}

func ExploreRange(start, end int64, next qp.Assemble) qp.Assemble {
	return qp.Map(1, exploreRange(start, end, next))
}

func exploreUnion(members []qp.Assemble) func(ipld.MapAssembler) {
	return func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreUnion, qp.List(int64(len(members)), func(la ipld.ListAssembler) {
			for _, member := range members {
				qp.ListEntry(la, member)
			}
		}))
	}
}

func BuildExploreUnion(np ipld.NodePrototype, members ...qp.Assemble) (ipld.Node, error) {
	return qp.BuildMap(np, 1, exploreUnion(members))
}

func ExploreUnion(members ...qp.Assemble) qp.Assemble {
	return qp.Map(1, exploreUnion(members))
}

func exploreFields(fields func(ma ipld.MapAssembler)) func(ma ipld.MapAssembler) {
	return func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreFields, qp.Map(1, func(ma ipld.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Fields, qp.Map(-1, fields))
		}))
	}
}

func BuildExploreFields(np ipld.NodePrototype, fields func(ma ipld.MapAssembler)) (ipld.Node, error) {
	return qp.BuildMap(np, 1, exploreFields(fields))
}

func ExploreFields(fields func(ma ipld.MapAssembler)) qp.Assemble {
	return qp.Map(1, exploreFields(fields))
}

func matcher(ma ipld.MapAssembler) {
	qp.MapEntry(ma, selector.SelectorKey_Matcher, qp.Map(0, func(ma ipld.MapAssembler) {}))
}

func BuildMatcher(np ipld.NodePrototype) (ipld.Node, error) {
	return qp.BuildMap(np, 1, matcher)
}

func Matcher() qp.Assemble {
	return qp.Map(1, matcher)
}
