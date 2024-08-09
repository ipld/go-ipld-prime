package basicnode_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/printer"
)

func TestList(t *testing.T) {
	tests.SpecTestListString(t, basicnode.Prototype.List)
}

func TestListAmendingBuilderNewNode(t *testing.T) {
	amender := basicnode.Prototype.List.AmendingBuilder(nil)

	err := amender.Append(basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Append(basicnode.NewString("dog"))
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Append(basicnode.NewString("eel"))
	if err != nil {
		t.Fatal(err)
	}
	listNode := amender.Build()
	expect := `list{
	0: string{"cat"}
	1: string{"dog"}
	2: string{"eel"}
}`
	actual := printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Update an element at the start
	err = amender.Set(0, basicnode.NewString("cow"))
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cow"}
	1: string{"dog"}
	2: string{"eel"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Update an element in the middle
	err = amender.Set(1, basicnode.NewString("fox"))
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cow"}
	1: string{"fox"}
	2: string{"eel"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Update an element at the end
	err = amender.Set(amender.Length(), basicnode.NewString("dog"))
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cow"}
	1: string{"fox"}
	2: string{"eel"}
	3: string{"dog"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Delete an element from the start
	err = amender.Remove(0)
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"fox"}
	1: string{"eel"}
	2: string{"dog"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Delete an element from the middle
	err = amender.Remove(1)
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"fox"}
	1: string{"dog"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Delete an element from the end
	err = amender.Remove(1)
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"fox"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Insert an element at the start
	err = amender.Insert(0, basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cat"}
	1: string{"fox"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Insert an element in the middle
	err = amender.Insert(1, basicnode.NewString("dog"))
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cat"}
	1: string{"dog"}
	2: string{"fox"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Insert an element at the end
	err = amender.Insert(amender.Length(), basicnode.NewString("eel"))
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cat"}
	1: string{"dog"}
	2: string{"fox"}
	3: string{"eel"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Access values of list, using index
	r, err := listNode.LookupByIndex(0)
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "cat", qt.Equals, must.String(r))

	// Access values of list, using PathSegment
	r, err = listNode.LookupBySegment(datamodel.ParsePathSegment("2"))
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "fox", qt.Equals, must.String(r))

	// Access updated value of list, using get
	r, err = amender.Get(3)
	if err != nil {
		t.Fatal(err)
	}
	qt.Check(t, "eel", qt.Equals, must.String(r))

	// Validate the node's prototype
	np := listNode.Prototype()
	qt.Check(t, fmt.Sprintf("%T", np), qt.Equals, "basicnode.Prototype__List")
}

func TestListAmendingBuilderExistingNode(t *testing.T) {
	listNode, err := qp.BuildList(basicnode.Prototype.List, -1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.String("fox"))
		qp.ListEntry(am, qp.String("cow"))
		qp.ListEntry(am, qp.String("deer"))
	})
	if err != nil {
		t.Fatal(err)
	}
	expect := `list{
	0: string{"fox"}
	1: string{"cow"}
	2: string{"deer"}
}`
	actual := printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)

	amender := basicnode.Prototype.List.AmendingBuilder(listNode)
	err = amender.Append(basicnode.NewString("cat"))
	if err != nil {
		t.Fatal(err)
	}
	newListNode := amender.Build()
	expect = `list{
	0: string{"fox"}
	1: string{"cow"}
	2: string{"deer"}
	3: string{"cat"}
}`
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	insertElems, err := qp.BuildList(basicnode.Prototype.List, -1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.String("eel"))
		qp.ListEntry(am, qp.String("dog"))
	})
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Insert(3, insertElems)
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"fox"}
	1: string{"cow"}
	2: string{"deer"}
	3: string{"eel"}
	4: string{"dog"}
	5: string{"cat"}
}`
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	err = amender.Append(basicnode.NewString("eel"))
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"fox"}
	1: string{"cow"}
	2: string{"deer"}
	3: string{"eel"}
	4: string{"dog"}
	5: string{"cat"}
	6: string{"eel"}
}`
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// The original node should not have been updated
	expect = `list{
	0: string{"fox"}
	1: string{"cow"}
	2: string{"deer"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)
}

func TestListAmendingBuilderCopiedNode(t *testing.T) {
	listNode, err := qp.BuildList(basicnode.Prototype.List, -1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.String("fox"))
		qp.ListEntry(am, qp.String("cow"))
		qp.ListEntry(am, qp.String("deer"))
	})
	if err != nil {
		t.Fatal(err)
	}
	amender := basicnode.Prototype.List.AmendingBuilder(listNode)
	newListNode := amender.Build()
	expect := `list{
	0: string{"fox"}
	1: string{"cow"}
	2: string{"deer"}
}`
	actual := printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	setElems, err := qp.BuildList(basicnode.Prototype.List, -1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.String("cat"))
		qp.ListEntry(am, qp.String("dog"))
	})
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Set(0, setElems)
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cat"}
	1: string{"dog"}
	2: string{"cow"}
	3: string{"deer"}
}`
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	insertElems, err := qp.BuildList(basicnode.Prototype.List, -1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.String("eel"))
		qp.ListEntry(am, qp.String("fox"))
	})
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Insert(1, insertElems)
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cat"}
	1: string{"eel"}
	2: string{"fox"}
	3: string{"dog"}
	4: string{"cow"}
	5: string{"deer"}
}`
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	appendElems, err := qp.BuildList(basicnode.Prototype.List, -1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.String("rat"))
	})
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Append(appendElems)
	if err != nil {
		t.Fatal(err)
	}
	expect = `list{
	0: string{"cat"}
	1: string{"eel"}
	2: string{"fox"}
	3: string{"dog"}
	4: string{"cow"}
	5: string{"deer"}
	6: string{"rat"}
}`
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Pass through an empty list. This should have no effect.
	appendElems, err = qp.BuildList(basicnode.Prototype.List, -1, func(am datamodel.ListAssembler) {})
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Append(appendElems)
	if err != nil {
		t.Fatal(err)
	}
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// Pass through a list containing another list
	nestedList, err := qp.BuildList(basicnode.Prototype.List, 1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.String("bat"))
	})
	if err != nil {
		t.Fatal(err)
	}
	appendElems, err = qp.BuildList(basicnode.Prototype.List, 1, func(am datamodel.ListAssembler) {
		qp.ListEntry(am, qp.Node(nestedList))
	})
	if err != nil {
		t.Fatal(err)
	}
	err = amender.Append(appendElems)
	if err != nil {
		t.Fatal(err)
	}
	// The new node should have been updated to have a list node at the end
	expect = `list{
	0: string{"cat"}
	1: string{"eel"}
	2: string{"fox"}
	3: string{"dog"}
	4: string{"cow"}
	5: string{"deer"}
	6: string{"rat"}
	7: list{
		0: string{"bat"}
	}
}`
	actual = printer.Sprint(newListNode)
	qt.Assert(t, actual, qt.Equals, expect)

	// The original node should not have been updated
	expect = `list{
	0: string{"fox"}
	1: string{"cow"}
	2: string{"deer"}
}`
	actual = printer.Sprint(listNode)
	qt.Assert(t, actual, qt.Equals, expect)
}
