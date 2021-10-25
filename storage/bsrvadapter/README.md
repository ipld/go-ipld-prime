bsrvadapter
===========

The `bsrvadapter` package/module is a small piece of glue code to connect
the `github.com/ipfs/go-blockservice` package, and packages implementing its interfaces,
forward into the `go-ipld-prime/storage` interfaces.

This can be used to rig systems like Bitswap up behind go-ipld-prime storage APIs.

(Whether or not this is a good idea is debatable.
It should be noted that both the `ipfs/go-blockservice` API,
as well as Bitswap in particular as an implementation,
are inherently prone to the infamous "N+1 Query Problem".
Treating a remote network fetch as equivalent to a local low-latency operation
just isn't a good idea for performance or predictability, no matter how you slice it.
Nonetheless: it's possible, using this code, if you really want to do it.)


Why structured like this?
-------------------------

See similar discussion in the `../dsadapter` module.


Which of `dsadapter` vs `bsadapter` vs `bsrvadapter` should I use?
------------------------------------------------------------------

None of them, ideally.
A direct implementation of the storage APIs will almost certainly be able to perform better than any of these adapters.

Failing that: use the adapter matching whatever you've got on hand in your code.

There is no correct choice.

dsadapter suffers avoidable excessive allocs in processing its key type,
due to choices in the interior of `github.com/ipfs/go-datastore`.
It is also unable to support streaming operation, should you desire it.

bsadapter and bsrvadapter both also suffer overhead due to their key type,
because they require a transformation back from the plain binary strings used in the storage API to the concrete go-cid type,
which spends some avoidable CPU time (and also, at present, causes avoidable allocs because of some interesting absenses in `go-cid`).
Additionally, they suffer avoidable allocs because they wrap the raw binary data in a "block" type,
which is an interface, and thus heap-escapes; and we need none of that in the storage APIs, and just return the raw data.
They are also unable to support streaming operation, should you desire it.

It's best to choose the shortest path and use the adapter to whatever layer you need to get to --
for example, if you really want to use a `go-datastore` implementation,
*don't* use `bsadapter` and have it wrap a `go-blockstore` that wraps a `go-datastore` if you can help it:
instead, use `dsadapter` and wrap the `go-datastore` without any extra layers of indirection.
You should prefer this because most of the notes above about avoidable allocs are true when
the legacy interfaces are communicating with each other, as well...
so the less you use the internal layering of the legacy interfaces, the better off you'll be.

Using a direct implementation of the storage APIs will suffer none of these overheads,
and so will always be your best bet if possible.

If you have to use one of these adapters, hopefully the performance overheads fall within an acceptable margin.
If not: we'll be overjoyed to accept help porting things.
