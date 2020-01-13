package ipld

// NodeAssembler is the interface that describes all the ways we can set values
// in a node that's under construction.
//
// To create a Node, you should start with a NodeBuilder (which contains a
// superset of the NodeAssembler methods, and can return the finished Node
// from its `Done` method).
//
// Why do both this and the NodeBuilder interface exist?
// When creating trees of nodes, recursion works over the NodeAssembler interface.
// This is important to efficient library internals, because avoiding the
// requirement to be able to return a Node at any random point in the process
// relieves internals from needing to implement 'freeze' features.
// (This is useful in turn because implementing those 'freeze' features in a
// language without first-class/compile-time support for them (as golang is)
// would tend to push complexity and costs to execution time; we'd rather not.)
type NodeAssembler interface {
	BeginMap(sizeHint int) (MapNodeAssembler, error)
	BeginList(sizeHint int) (ListNodeAssembler, error)
	AssignNull() error
	AssignBool(bool) error
	AssignInt(int) error
	AssignFloat(float64) error
	AssignString(string) error
	AssignBytes([]byte) error

	Assign(Node) error // if you already have a completely constructed subtree, this method puts the whole thing in place at once.

	Style() NodeStyle // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
}

type MapNodeAssembler interface {
	Insert(string, Node) error
	InsertComplexKey(Node, Node) error
	AssembleInsertion() (NodeAssembler, NodeAssembler) // NOT POSSIBLE: latter may depend on actions on former.

	Done() error

	KeyStyle() NodeStyle   // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
	ValueStyle() NodeStyle // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
}

type ListNodeAssembler interface {
	Append(Node) error
	AssembleValue() NodeAssembler

	Done() error

	ValueStyle() NodeStyle // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
}

type NodeBuilder interface {
	NodeAssembler
	Build() (Node, error)

	// Resets the builder.  It can hereafter be used again.
	// Reusing a NodeBuilder can reduce allocations and improve performance.
	// If you're not reusing the builder, don't call this.
	//
	// (Authors of generic algorithms handling maps: if your map has complex
	// keys (e.g., you need to use a NodeBuilder and the InsertComplexKey
	// function rather than just Insert), it's often particularly impact
	// useful to
	// no
	// avoiding only one of the two allocs per key is Not Good
	// okay... i can't think of a way to generically copy map keys without allocs.
	// unless all maps basically contain the key twice: once in the map and once in a (second) slice.
	// or if we just... completely reimplement maps.  but let's not.
	// i guess this is also something we can/should defer...
	// key takeaways for now are:
	// - yes, you might actually need a values-only iterator, because holy crap
	//   - might be worth an experiment to see if unused things can be escape-analyzed out; but i'm betting against.
	//     - we can just... deliver the thing ignoring this, and do this later.  purely additive future work, no backtracking.
	// - yes, you definitely need a way to stream in parts of a complex key... a builder outside and InsertComplexKey is not great.
	// - optimizing for the key-comes-from-another-map mode might be trickier than the rest... but that's mostly on the source side rather than accept side.
	Reset()
}

// sidequest for ergonomics test, even though this is priority 2 for these APIs
func demo() {
	mustChill := func(error) {}
	var nb NodeBuilder
	func(ma MapNodeAssembler) {
		ka, va := ma.AssembleInsertion()
		mustChill(ka.AssignString("key"))
		func(ma MapNodeAssembler) {
			ka, va := ma.AssembleInsertion()
			mustChill(ka.AssignString("nested"))
			mustChill(va.AssignBool(true))
			ma.Done()
		}(va.BeginMap())
		ka, va = ma.AssembleInsertion() // calling this repeatedly will be annoying.  (for `:=`/`=` reasons, mainly.  also: disrupts lhs flow.)
		mustChill(ka.AssignString("secondkey"))
		mustChill(va.AssignString("morevalue"))
		ma.Done()
	}(nb.BeginMap())
	result, err := nb.Build()
	_, _ = result, err
	// basically no part of the above works if you mentally replace `mustChill` with the three-liner that needs to return.
	// so this whole demo is wishful thinking that's not particularly plausible;
	// if we want indentation like this, it's gonna need to come from wrappers (like the already explored 'fluent' package).
}

type AltMapNodeAssembler interface {
	Insert(string, Node) error
	InsertComplexKey(Node, Node) error
	AssembleKey() NodeAssembler   // must be followed by call to AssembleValue.
	AssembleValue() NodeAssembler // must be called after AssembleKey.
	// ^ this is attempting to fix the "lhs flow" issue, but doesn't touch the big compile error above around BeginMap returning error.
	// ^ also probably require a word or two less memory to implement (i think).

	Done() error

	KeyStyle() NodeStyle   // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
	ValueStyle() NodeStyle // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
}

// Complex keys: What do they come from?  (What arrre they _good_ for? (Absolutely nothin, say it again))
// - in the Data Model, they don't exist.  They just... don't exist.
//   - json and javascript could never real deal reliably with numbers, so we just... nope.
//   - maps keyed by multiple kinds are also beyond the pale in many host languages, so, again... nope.
// - in the schema types layer, they can exist.
//   - a couple things can reach it:
//     - `type X struct {...} representation stringjoin`
//     - `type X struct {...} representation stringpairs`
//     - `type X map {...} representation stringpairs` // *maybe*.  go won't allow using this as a map key except in string form anyway.
//     - we don't know what the syntax would look like for a type-level int, but, haven't ruled it out either.
//   - but when feeding data in via representation: it's all strings, of course.
//   - if we have codegen and are app author, we can use native methods that pass-by-value.
//   - so it's ONLY when doing generic code, or typed but not using codegen, that we face these apis.
// - and it's ONLY -- this should go without saying, but let's say it -- relevant when thinking about map keys.
//   - structs can't have complex keys.  field names are strings.  done.
//   - lists can't have complex keys.  obvious category error.
//   - enums can't have complex keys.  because we said so; or, obvious category error; take your pick.
//   - why is this important?  well, it means complex keys only exist in a place where values are all going to use the same style/builder.
//     - which means `AssembleInsertion() (NodeAssembler, NodeAssembler)` is actually kinda back on the table.
//       - though since it would only work in *some* cases... still serious questions about whether that'd be a good API to show off.
//       - we still also need to be able to get a NodeAssembler for streaming in value content when handling structs, so, then we'd end up with two APIs for this?
//         - yeah, this design still does not go good places; abort.
