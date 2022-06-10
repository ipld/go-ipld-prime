package amend

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal/patch"
)

var addTests = []struct {
	size int
	num  int
}{
	{size: 100, num: 1},
	{size: 100, num: 10},
	{size: 100, num: 100},
	{size: 1000, num: 10},
	{size: 1000, num: 100},
	{size: 1000, num: 1000},
	{size: 10000, num: 100},
	{size: 10000, num: 1000},
	{size: 10000, num: 10000},
}

var removeTests = []struct {
	size int
	num  int
}{
	{size: 100, num: 1},
	{size: 100, num: 10},
	{size: 100, num: 100},
	{size: 1000, num: 10},
	{size: 1000, num: 100},
	{size: 1000, num: 1000},
	{size: 10000, num: 100},
	{size: 10000, num: 1000},
	{size: 10000, num: 10000},
}

var replaceTests = []struct {
	size int
	num  int
}{
	{size: 100, num: 1},
	{size: 100, num: 10},
	{size: 100, num: 100},
	{size: 1000, num: 10},
	{size: 1000, num: 100},
	{size: 1000, num: 1000},
	{size: 10000, num: 100},
	{size: 10000, num: 1000},
	{size: 10000, num: 10000},
}

func BenchmarkAmend_Map_Add(b *testing.B) {
	for _, v := range addTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildMap(basicnode.Prototype.Any, int64(v.size), func(ma datamodel.MapAssembler) {
				for i := 0; i < v.size; i++ {
					qp.MapEntry(ma, "key-"+strconv.Itoa(i), qp.String("value-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				a := NewAmender(tn)
				for i := 0; i < v.num; i++ {
					err = EvalOne(a, patch.Operation{
						Op:    patch.Op_Add,
						Path:  datamodel.ParsePath("/new-key-" + strconv.Itoa(i)),
						Value: basicnode.NewString("new-value-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("amend did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(a.Build(), dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkPatch_Map_Add(b *testing.B) {
	for _, v := range addTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildMap(basicnode.Prototype.Any, int64(v.size), func(ma datamodel.MapAssembler) {
				for i := 0; i < v.size; i++ {
					qp.MapEntry(ma, "key-"+strconv.Itoa(i), qp.String("value-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				for i := 0; i < v.num; i++ {
					tn, err = patch.EvalOne(tn, patch.Operation{
						Op:    patch.Op_Add,
						Path:  datamodel.ParsePath("/new-key-" + strconv.Itoa(i)),
						Value: basicnode.NewString("new-value-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("patch did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(tn, dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkAmend_List_Add(b *testing.B) {
	for _, v := range addTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildList(basicnode.Prototype.Any, int64(v.size), func(la datamodel.ListAssembler) {
				for i := 0; i < v.size; i++ {
					qp.ListEntry(la, qp.String("entry-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				a := NewAmender(tn)
				for i := 0; i < v.num; i++ {
					err = EvalOne(a, patch.Operation{
						Op:    patch.Op_Add,
						Path:  datamodel.ParsePath("/0"), // insert at the start for worst-case
						Value: basicnode.NewString("new-entry-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("amend did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(a.Build(), dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkPatch_List_Add(b *testing.B) {
	for _, v := range addTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildList(basicnode.Prototype.Any, int64(v.size), func(la datamodel.ListAssembler) {
				for i := 0; i < v.size; i++ {
					qp.ListEntry(la, qp.String("entry-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				for i := 0; i < v.num; i++ {
					tn, err = patch.EvalOne(tn, patch.Operation{
						Op:    patch.Op_Add,
						Path:  datamodel.ParsePath("/0"), // insert at the start for worst-case
						Value: basicnode.NewString("new-entry-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("patch did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(tn, dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkAmend_Map_Remove(b *testing.B) {
	for _, v := range removeTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildMap(basicnode.Prototype.Any, int64(v.size), func(ma datamodel.MapAssembler) {
				for i := 0; i < v.size; i++ {
					qp.MapEntry(ma, "key-"+strconv.Itoa(i), qp.String("value-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				a := NewAmender(tn)
				for i := 0; i < v.num; i++ {
					err = EvalOne(a, patch.Operation{
						Op:   patch.Op_Remove,
						Path: datamodel.ParsePath("/key-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("amend did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(a.Build(), dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkPatch_Map_Remove(b *testing.B) {
	for _, v := range removeTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildMap(basicnode.Prototype.Any, int64(v.size), func(ma datamodel.MapAssembler) {
				for i := 0; i < v.size; i++ {
					qp.MapEntry(ma, "key-"+strconv.Itoa(i), qp.String("value-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				for i := 0; i < v.num; i++ {
					tn, err = patch.EvalOne(tn, patch.Operation{
						Op:   patch.Op_Remove,
						Path: datamodel.ParsePath("/key-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("patch did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(tn, dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkAmend_List_Remove(b *testing.B) {
	for _, v := range removeTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildList(basicnode.Prototype.Any, int64(v.size), func(la datamodel.ListAssembler) {
				for i := 0; i < v.size; i++ {
					qp.ListEntry(la, qp.String("entry-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				a := NewAmender(tn)
				for i := 0; i < v.num; i++ {
					err = EvalOne(a, patch.Operation{
						Op:   patch.Op_Remove,
						Path: datamodel.ParsePath("/0"), // remove from the start for worst-case
					})
					if err != nil {
						b.Fatalf("amend did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(a.Build(), dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

// TODO: Investigate panic
//func BenchmarkPatch_List_Remove(b *testing.B) {
//	for _, v := range removeTests {
//		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
//			n, _ := qp.BuildList(basicnode.Prototype.Any, int64(v.size), func(la datamodel.ListAssembler) {
//				for i := 0; i < v.size; i++ {
//					qp.ListEntry(la, qp.String("entry-"+strconv.Itoa(i)))
//				}
//			})
//			var err error
//			for r := 0; r < b.N; r++ {
//				tn := n
//				for i := 0; i < v.num; i++ {
//					tn, err = patch.EvalOne(tn, patch.Operation{
//						Op:    patch.Op_Remove,
//						Path:  datamodel.ParsePath("/0"), // remove from the start for worst-case
//					})
//					if err != nil {
//						b.Fatalf("patch did not apply: %s", err)
//					}
//				}
//				output, err := ipld.Encode(tn, dagjson.EncodeOptions{
//					EncodeLinks: true,
//					EncodeBytes: true,
//					MapSortMode: codec.MapSortMode_None,
//				}.Encode)
//				log.Printf("json: %s", output)
//				if err != nil {
//					b.Errorf("failed to serialize result: %s", err)
//				}
//			}
//		})
//	}
//}

func BenchmarkAmend_Map_Replace(b *testing.B) {
	for _, v := range replaceTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildMap(basicnode.Prototype.Any, int64(v.size), func(ma datamodel.MapAssembler) {
				for i := 0; i < v.size; i++ {
					qp.MapEntry(ma, "key-"+strconv.Itoa(i), qp.String("value-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				a := NewAmender(tn)
				for i := 0; i < v.num; i++ {
					err = EvalOne(a, patch.Operation{
						Op:    patch.Op_Replace,
						Path:  datamodel.ParsePath("/key-" + strconv.Itoa(rand.Intn(v.size))),
						Value: basicnode.NewString("new-value-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("amend did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(a.Build(), dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkPatch_Map_Replace(b *testing.B) {
	for _, v := range replaceTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildMap(basicnode.Prototype.Any, int64(v.size), func(ma datamodel.MapAssembler) {
				for i := 0; i < v.size; i++ {
					qp.MapEntry(ma, "key-"+strconv.Itoa(i), qp.String("value-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				for i := 0; i < v.num; i++ {
					tn, err = patch.EvalOne(tn, patch.Operation{
						Op:    patch.Op_Replace,
						Path:  datamodel.ParsePath("/key-" + strconv.Itoa(rand.Intn(v.size))),
						Value: basicnode.NewString("new-value-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("patch did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(tn, dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkAmend_List_Replace(b *testing.B) {
	for _, v := range replaceTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildList(basicnode.Prototype.Any, int64(v.size), func(la datamodel.ListAssembler) {
				for i := 0; i < v.size; i++ {
					qp.ListEntry(la, qp.String("entry-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				a := NewAmender(tn)
				for i := 0; i < v.num; i++ {
					err = EvalOne(a, patch.Operation{
						Op:    patch.Op_Replace,
						Path:  datamodel.ParsePath("/" + strconv.Itoa(rand.Intn(v.size))),
						Value: basicnode.NewString("new-entry-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("amend did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(a.Build(), dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}

func BenchmarkPatch_List_Replace(b *testing.B) {
	for _, v := range replaceTests {
		b.Run(fmt.Sprintf("inputs: %v", v), func(b *testing.B) {
			n, _ := qp.BuildList(basicnode.Prototype.Any, int64(v.size), func(la datamodel.ListAssembler) {
				for i := 0; i < v.size; i++ {
					qp.ListEntry(la, qp.String("entry-"+strconv.Itoa(i)))
				}
			})
			var err error
			for r := 0; r < b.N; r++ {
				tn := n
				for i := 0; i < v.num; i++ {
					tn, err = patch.EvalOne(tn, patch.Operation{
						Op:    patch.Op_Replace,
						Path:  datamodel.ParsePath("/" + strconv.Itoa(rand.Intn(v.size))),
						Value: basicnode.NewString("new-entry-" + strconv.Itoa(i)),
					})
					if err != nil {
						b.Fatalf("patch did not apply: %s", err)
					}
				}
				_, err = ipld.Encode(tn, dagjson.EncodeOptions{
					EncodeLinks: true,
					EncodeBytes: true,
					MapSortMode: codec.MapSortMode_None,
				}.Encode)
				if err != nil {
					b.Errorf("failed to serialize result: %s", err)
				}
			}
		})
	}
}
