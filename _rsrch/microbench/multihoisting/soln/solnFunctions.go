package solution

type Node interface {
	LookupByString(key string) Node
}

type NodeBuilder interface {
	InsertByString(key string, value Node)
	Build() Node
}

func (n *Stroct) LookupByString(key string) Node {
	switch key {
	case "foo":
		return &n.foo
	case "bar":
		return &n.bar
	default:
		panic("no")
	}
}

func (n *Strooct) LookupByString(key string) Node {
	panic("nyi")
}

func (n *String) LookupByString(key string) Node {
	panic("nyi")
}

func NewStroctBuilder() NodeBuilder {
	r := _Stroct__Racker{}
	return &r._Stroct__Builder
}

func (nb *_Stroct__Builder) InsertByString(key string, value Node) {
	switch key {
	case "foo":
		if nb.isset_foo {
			panic("cannot set field repeatedly!") // surprisingly, could make this optional without breaking memviz semantics.
		}
		if &nb.d.foo != value.(*Strooct) { // shortcut: if shmem, no need to memcopy self to self!
			nb.d.foo = *value.(*Strooct) // REVIEW: maybe this isn't necessary here and only should be checked in the racker?
		}
		// interestingly, no need to set br.frz_foo... nb.b isn't revealed yet, so, while
		nb.isset_foo = true
	case "bar":
		nb.d.bar = *value.(*String)
		nb.isset_bar = true
	default:
		panic("no")
	}
}
func (br *_Stroct__Racker) InsertByString(key string, value Node) {
	switch key {
	case "foo":
		if &br.d.foo == value.(*Strooct) { // not just a shortcut: if shmem, this branch insert is ONLY possible WHEN the field is frozen.
			br.isset_foo = true // FIXME straighten this out to dupe less plz?!  also should panic if already set, consistently (or not).
		}
		if br.frz_foo {
			panic("cannot set field which has been frozen due to shared immutable memory!")
		}
		// TODO finish and sanity check in morning.
		// ... seems like there's a lot less shared code than i thought.
		// ... resultingly, i'm seriously questioning if builders and rackers deserve separate types.
		//   ... what's that for? reducing the (fairly temporary) size of builders if used sans racking style?  does that... matter?
		// no, by the morning light i think this can in fact just be simplified a lot and code share will emerge.
	default:
		panic("no")
	}
}

func (nb *_Stroct__Builder) Build() Node {
	return nb.d
}
func (br *_Stroct__Racker) Build() Node {
	// TODO freeze something.  but what?  do i really need a pointer "up" that to be fed into here?
	// is that even a pointer?  if it depends on what this value is contained in... remember, for maps and lists, it's different than structs: bools ain't it.
	// if we try to use a function pointer here: A) you're kidding.. that will not be fast.. B) check to see if it causes fucking allocs, because it fucking might, something about grabbing methods causing closures.

	// AH HA.  Obrigado!:
	//  We can actually use the invalidation of the 'd' pointer as our signal, and do the freezing update the next time either Set or GetBuilderForValue is called.

	// For structs, since we'll keep a builder/racker per field, that's enough state already.  You'll be able to get the childbuilders for multiple fields simultaneously (though you shouldn't rely on it, in general).
	// For lists, there will be a frozen offset and an isset offset.  The former will be able to advance well beyond the latter.  You'll only be allowed to have one childbuilder at a time.
	// For maps, we'll still do copies and allocs (reshuffling maps with large embeds is its own cost center usually better avoided).  You'll only be able to use one keybuilder at a time (keys will be copied when done), and one value builder (but each use of the value builder will incur an alloc for its 'd' field).
	//  And here we hit one more piece of Fun that might require a small break to current interfaces: in order to make it possible for keys to reuse a swap space and not get shifted to heap... we can't return them.
	//   Which means no ChildBuilderForKey method at all: that interface doesn't fly.
	//    Oh dear.  This is unpleasant.  More thought needed on this.
	//     ... Hang on.  Have I imagined this more complicated than it is?
	//      For structs, the key is always a string.
	//      For maps with enum keys... always a string.  (...I don't know how to say enums of int repr would work; I don't think they do, Because Data Model rules.)
	//      For maps with union keys... okay, we're gonna require that those act like struct keys (unions used in keys will have to have nonptr internals; a fun detail).
	//      For maps with struct keys... okay, sure it has to be a string from repr land... but also yes, we have to accept the mappy form, for typed level copies and transforms.
	//       Yup.  That last bit means a need for swap space while assembling it.  Which we don't want to be forced to return as a node, because that would cause a heap alloc.  Which we'd immediately undo since map keys need to be by value to behave correctly.  Ow.
	//        I... don't know what to do about this.  We could have the Build method flat out return a nil, and just "document" that.  But it's pretty unappealing to make such an edge case.

	// Also, yeah, we really do need an InsertByString method on MapBuilder.  Incurring some nonsense boxing for string keys in structs is laughable.
	//  If you're thinking a workaroudn such as having a single swap space for building a single justString key for temporary use would help... no, sadly: "single swap space" plus visibility model won't jive like that.
	//   (And even if it did, ever asking an end user to write that much boilerplate is still pretty crass... as well as easily avoidable for minimal library code size cost.)
	return br.d
}
