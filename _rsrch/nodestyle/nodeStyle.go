package nodestyle

import (
	"context"
	"io"
)

type ReprKind uint8

type Node interface {
	ReprKind() ReprKind
	LookupString(key string) (Node, error)
	Lookup(key Node) (Node, error)
	LookupIndex(idx int) (Node, error)
	Length() int
	IsUndefined() bool
	IsNull() bool
	AsBool() (bool, error)
	AsInt() (int, error)
	AsFloat() (float64, error)
	AsString() (string, error)
	AsBytes() ([]byte, error)
	AsLink() (Link, error)

	Prototype() NodePrototype // note!  replaces `NodeBuilder` method!
}

// Prototype is the information that lets you make more of them, and do some basic behavioral inspection;
// Type is the information that lets you understand how a group of nodes is related in a schema and rules it will follow.
type NodePrototype interface {
	NewBuilder() NodeBuilder // allocs!  (probably.  sometimes the alloc is still later.)
}

type NodePrototypeSupportingAmend interface {
	AmendingBuilder(base Node) NodeBuilder
}

// all the error methods here are still a serious question.  i don't really like them.  you can 'must' them away, but... squick?
//    ... we should make a muster that stores them somewhere rather than panics.  or takes a type parameter for its panic, or `func(error) boxedError`.  or something.
type NodeBuilder interface {
	BeginMap() (MapBuilder, error)   // note the "amend" options are gone -- now do it with feature detection on NodePrototypeSupportingAmend, instead!
	BeginList() (ListBuilder, error) // note the "amend" options are gone -- now do it with feature detection on NodePrototypeSupportingAmend, instead!
	CreateNull() (Node, error)
	CreateBool(bool) (Node, error)
	CreateInt(int) (Node, error)
	CreateFloat(float64) (Node, error)
	CreateString(string) (Node, error)
	CreateBytes([]byte) (Node, error)
	CreateLink(Link) (Node, error) // fixme this is dumb and all links should already be nodes; either that or their creation should hinge here, rather than elsewhere and be awkwardly doubled.

	Prototype() NodePrototype // it's unlikely this will often be needed, i think, but it's here nonetheless.  (maybe generic transform will find this easier to use than the one on the node?  i think both should end up in reach on the stack, but not sure.)
}

type MapBuilder interface {
	// in question whether this needs error returns and/or possibly a 'Done' method at all
	// philosophically: is it worse if this cursor make cause an error when a method is called on its parent object (vs if we just put all cursor behavior in one object flatly)?

	// how maps do key recursion is one major question;
	// and whether the 'Done' method returns a Node is the other.
	// the nodebuilder above seems to think things return Node, but that doesn't have to be.
	// size hinting also got forgotten again.

	// problem with trying to put the results elsewhere is... where?
	//  the NodeBuilder above *doesn't have* a bind operation or Done method that could serve.
	//   didn't we have a draft somewhere that... did?

	// yeah, we did.  that's the 'nodeassembler' draft.
	// currently thinking: best hybrid is:
	// - yes, nodeStyle
	// - yes, nodeAssembler vs nodeBuilder -- BUT:
	// - nodeAssembler still return error promptly all over the place.
	// the reasons these are the preferable choices are:
	// - nodeStyle still dekinks about a dozen things
	// - not returning intermediate nodes simplifies the builder/freezing cycle radically (basically, removes it) and this means going faster and with less code.
	// - returning errors promptly all over the place reduces intermediate assembler structs needing to store errors just to curry weirdly and temporarily (esp. on scalar/leaf node construction).
	// - and yes, in most generic algorithms, we're going to want to return errors completely promptly -- the position in the other side of a hylomorphism is relevant.
	// - these interfaces are mostly for high performance in generic algorithms.  we can attempt to build additional friendliness for application programmers on at higher or alternative laters.

}

type ListBuilder interface {
}

type Link interface {
	// not exploring it here, but `LinkBuilder` might also make more sense if renamed to `LinkPrototype`, and split into
	// but there's a lot of things about that interface that still don't "feel" friendly (even if they're correct).
	// ... okay, exploring it here.
	// ... it's possible that traversal.Config.LinkNodeBuilderChooser is actually still right; you don't need or want to see the NodePrototype there?
	//   unless we want to use the NodePrototype has the hub for feature detection for other fastpaths?  which we probably... indeed might want.  hm.

	Prototype() LinkPrototype
}

type LinkPrototype interface {
	// unsure what goes here, tbh.
}

type StorageLoader func(ctx context.Context, lnk Link, lnkCtx LinkContext) (io.Reader, error) // just handles the concept of bytes -- might have an internal `(Link)->(filepath.Path)` func, but `Link` is otherwise opaque to it.

type StorageWriter func(ctx context.Context, lnkCtx LinkContext) (io.Writer, StorageCommitter, error)
type StorageCommitter func(Link) error

type LinkContext struct {
	LinkPath   Path
	LinkNode   Node // has the Link again, but also might have type info // always zero for writing new nodes, for obvi reasons. // dubious if this is needed; would rather make Link impls just also *be* Node.
	ParentNode Node
}

type Path string
