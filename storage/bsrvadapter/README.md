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

See `../README_adapters.md` for details about why adapter code is needed,
why this is in a module, why it's here, etc.


Which of `dsadapter` vs `bsadapter` vs `bsrvadapter` should I use?
------------------------------------------------------------------

In short: you should prefer direct implementations of the storage APIs
over any of these adapters, if one is available with the features you need.

Otherwise, if that's not an option (yet) for some reason,
use whichever adapter gets you most directly connected to the code you need.

See `../README_adapters.md` for more details and discussion.
