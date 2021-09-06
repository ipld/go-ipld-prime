/*
	The 'node' package gathers various general purpose Node implementations;
	the first one you should jump to is 'node/basicnode'.

	There's no code in this package itself; it's just for grouping.

	The `Node` interface itself is in the `ipld` package,
	which is the parent of this.

	The 'node/mixins' package contains reusable component code for building
	your own node implementations, should you desire to do so.
	This includes standardized behavioral tests (!), which are
	in the 'node/mixins/tests' package.

	Other planned subpackages include:
	a cbor-native Node implementation (which can optimize performance in some
	cases by lazily parsing serial	data, and also retaining it as byte slice
	references for minimizing reserialization work for small mutations);
	a Node implementation which works over golang native types by use of reflection;
	a Node implementation which supports Schema type constraints and works
	without compile-time/codegen support by delegating storage to another Node implementation;
	etc.

	You can create your own Node implementations, too.
	There's nothing special about being in this package.

	Other Node implementations not found here may include those which
	are output from Schema-powered codegen!
*/
package node
