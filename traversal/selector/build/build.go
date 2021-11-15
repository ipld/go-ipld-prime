package build

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	selector "github.com/ipld/go-ipld-prime/traversal/selector"
)

type BuildFn func(na datamodel.NodeAssembler)

func ExploreRecursiveEdge() BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreRecursiveEdge, qp.Map(0, func(datamodel.MapAssembler) {}))
	})
}

func ExploreRecursive(limit selector.RecursionLimit, sequence BuildFn) BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreRecursive, qp.Map(2, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Limit, qp.Map(1, func(ma datamodel.MapAssembler) {
				switch limit.Mode() {
				case selector.RecursionLimit_Depth:
					qp.MapEntry(ma, selector.SelectorKey_LimitDepth, qp.Int(limit.Depth()))
				case selector.RecursionLimit_None:
					qp.MapEntry(ma, selector.SelectorKey_LimitNone, qp.Map(0, func(na datamodel.MapAssembler) {}))
				default:
					panic("Unsupported recursion limit type")
				}
			}))
			qp.MapEntry(ma, selector.SelectorKey_Sequence, sequence)
		}))
	})
}

func ExploreAll(next BuildFn) BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreAll, qp.Map(1, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Next, next)
		}))
	})
}

func ExploreIndex(index int64, next BuildFn) BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreIndex, qp.Map(2, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Index, qp.Int(index))
			qp.MapEntry(ma, selector.SelectorKey_Next, next)
		}))
	})
}

func ExploreRange(start, end int64, next BuildFn) BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreRange, qp.Map(3, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Start, qp.Int(start))
			qp.MapEntry(ma, selector.SelectorKey_End, qp.Int(end))
			qp.MapEntry(ma, selector.SelectorKey_Next, next)
		}))
	})
}

func ExploreUnion(members ...BuildFn) BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreUnion, qp.List(int64(len(members)), func(la datamodel.ListAssembler) {
			for _, member := range members {
				qp.ListEntry(la, member)
			}
		}))
	})
}

func ExploreFields(fields map[string]BuildFn) BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_ExploreFields, qp.Map(1, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, selector.SelectorKey_Fields, qp.Map(int64(len(fields)), func(ma datamodel.MapAssembler) {
				for field, selector := range fields {
					qp.MapEntry(ma, field, selector)
				}
			}))
		}))
	})
}

func Matcher() BuildFn {
	return qp.Map(1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, selector.SelectorKey_Matcher, qp.Map(0, func(ma datamodel.MapAssembler) {}))
	})
}

func Path(path datamodel.Path, next BuildFn) BuildFn {
	if path.Len() == 0 {
		return next
	}
	segment, remaining := path.Shift()
	index, err := segment.Index()
	if err != nil {
		return ExploreIndex(index, Path(remaining, next))
	} else {
		return ExploreFields(map[string]BuildFn{
			segment.String(): Path(remaining, next),
		})
	}
}

func Build(buildFn BuildFn) (_ datamodel.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok {
				err = rerr
			} else {
				// A reasonable fallback, for e.g. strings.
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	nb := basicnode.Prototype.Map.NewBuilder()
	buildFn(nb)
	return nb.Build(), nil
}
