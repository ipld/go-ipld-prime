package ipldfree

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/tests"
)

func Test(t *testing.T) {
	tests.TestNodes(t, func() ipld.MutableNode {
		return &Node{}
	})
}
