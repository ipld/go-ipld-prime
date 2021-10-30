bsadapter
=========

The `bsadapter` package/module is a small piece of glue code to connect
the `github.com/ipfs/go-blockstore` package, and packages implementing its interfaces,
forward into the `go-ipld-prime/storage` interfaces.


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
