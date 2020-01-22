package solution

// -- ipld schema -->
/*
	type Stroct struct {
		foo Strooct
		bar String
	}

	type Strooct struct {
		zot  String
		zam  String
		zems Strems
		zigs Zigs
		zee  Zahn
	}

	type Strems [String]

	type Zigs {String:Zahn}

	type Zahn struct {
		bahn String
	}
*/

// -- the readable types -->

type Stroct struct {
	foo Strooct
	bar String
}
type Strooct struct {
	zot  String
	zam  String
	zems Strems
	zigs Zigs
	zee  Zahn
}
type String struct {
	x string
}
type Strems struct {
	x []String
}
type Zigs struct {
	x map[String]Zahn
}
type Zahn struct {
	bahn String
}

// -- the builders alone -->

type _Stroct__Builder struct {
	d *Stroct // this pointer aims into the thing we're building (it's as yet unrevealed).  it will be nil'd when we reveal it.

	isset_foo bool
	isset_bar bool
}
type _Strooct__Builder struct {
	d *Strooct

	isset_zot  bool
	isset_zam  bool
	isset_zems bool
	isset_zigs bool
	isset_zee  bool
}
type _Strems__Builder struct {
	d *Strems
	// TODO
}
type _String__Builder struct {
	// okay, this one is a gimme: data contains only a ptr itself, effectively.
	// TODO: still might need a 'd' pointer, in case you're assigning into a list and wanna save boxing allocs?
	//  ... so long as we're doing wrapper types (to block blind casting), we're gonna have those boxing alloc concerns.
}
type _Zigs__Builder struct {
	d *Zigs
	// TODO
}

// -- the rackerized builders -->

type _Stroct__Racker struct {
	_Stroct__Builder // most methods come from this, but child-builder getters will be overriden.

	cb_foo  _Strooct__Racker // provides child builder for field 'foo'.
	frz_foo bool             // if true, must never yield cb_foo again.  becomes true on cb_foo.Build *or* assignment to field (latter case reachable if the value was made without going through cb_foo).
}
type _Strooct__Racker struct {
	_Strooct__Builder // most methods come from this, but child-builder getters will be overriden.

	// TODO: we might still actually need builders for scalars.  if it's got a wrapper struct, it would incur boxing.  damnit.
	zems _Strems__Racker
	zigs _Zigs__Racker
	zee  _Zahn__Racker
}
type _Strems__Racker struct {
	// TODO didn't finish
}
type _Zigs__Racker struct {
	// TODO didn't finish
}
type _Zahn__Racker struct {
	// TODO didn't finish
}

// right, here's one wild ride we haven't addressed yet:
// if you build a thing that resides in racker-operated memory, you get a node.
// so far so good, and you can even use it multiple places.
// if assignments go through a Maybe struct?
// ... actually, this is all fine.
// the 'MaybeFoo' structs should store pointers to the thing.  done.
// if the thing originated in racker-operated memory, this is free;
// if it didn't, it's a cost you would've hit somewhere else already anyway too.
// done.  it's fine.
