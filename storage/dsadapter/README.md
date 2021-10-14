dsadapter
=========

The `dsadapter` package/module is a small piece of glue code to connect
the `github.com/ipfs/go-datastore` package, and packages implementing its interfaces,
forward into the `go-ipld-prime/storage` interfaces.

For example, this can be used to use "flatfs" and other datastore plugins
with go-ipld-prime storage APIs.

Why structured like this?
-------------------------

Why are there layers of interface code?
The `go-ipld-prime/storage` interfaces are a newer generation,
and improves on several things vs `go-datastore`.  (See other docs for that.)

Why is this code in a shared place?
The glue code to connect `go-datastore` to the new `go-ipld-prime/storage` APIs
is fairly minimal, but there's also no reason for anyone to write it twice,
so we want to put it somewhere easy to share.

Why does this code has its own go module?
A separate module is used because it's important that go-ipld-prime can be used
without forming a dependency on `go-datastore`.
(We want this so that there's a reasonable deprecation pathway -- it must be
possible to write new code that doesn't take on transitive dependencies to old code.)

Why does this code exist here, in this git repo?
We put this separate module in the same git repo as `go-ipld-prime`... because we can.
Technically, neither this module nor the go-ipld-prime module depend on each other --
they just have interfaces that are aligned with each other -- so it's very easy to
hold them as separate go modules in the same repo, even though that can otherwise sometimes be tricky.
