package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"

	. "github.com/warpfork/go-wish"
)

type MutableNodeFactory func() ipld.MutableNode

func TestNodes(t *testing.T, nfac MutableNodeFactory) {
	n1 := nfac()
	n1.SetString("asdf")
	Wish(t, fluent.WrapNode(n1).AsString(), ShouldEqual, "asdf")
}
