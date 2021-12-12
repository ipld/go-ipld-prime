benchmarks
==========

This is a small module that pulls in a bunch of storage implementations,
as well as legacy implementations via the adapter modules,
and benchmarks all of them on the same benchmarks.

There's no reason to import this code,
so the go.mod file uses relative paths shamelessly.
(You can create your own benchmarks using the code in `../tests`,
which contains most of the engine; this package is just tables of setup.)


What variations do the benchmarks exercise?
------------------------------------------

- the various storage implementations!
	- in some cases: variations of parameters to individual storage implementations.  (TODO)
- puts and gets.  (TODO: currently only puts.)
- various distributions of data size.  (TODO)
- block mode vs streaming mode.  (TODO)
- end-to-end use via linksystem with small cbor objects.  (TODO)
	- (this measures a lot of things that aren't to do with the storage itself -- but is useful to contextualize things.)

Running the benchmarks on variations in hardware and filesystem may also be important!
Many of these storage systems use the disk in some way.


Why is the module structured like this?
----------------------------------------

Because many of the storage implementations are also their own modules,
and we don't want to have the go-ipld-prime module pull in a huge universe of transitive dependencies.

See similar discussion in `../README_adapters.md`.

It may be worth pulling this out into a new git repo in the future,
especially if we want to add more and more implementations to what we benchmark,
or develop additional tools for deploying the benchmark on varying hardware, etc.
For now, it incubates here.
