## v0.0.6 / 2017-03-23

* Improve error handling during encryption and decryption.
* Fix security regression in new BoltDB backend.

## v0.0.5 / 2017-03-22

* Make BoltDB the default.

## v0.0.4 / 2017-03-21

* Fixed styling of "wrong password" modal.
* Remove password and close buttons from tab index to optimize on the access pattern.
* Start publishing a docker container.

## v0.0.3 / 2017-03-20

This release adds support for a new
[BoltDB](https://github.com/boltdb/bolt) backend.  This backend was
*largely* added in order enable support for operating systems other
than Linux. The original backend used by Notable was
[Sqlite3](https://www.sqlite.org/) which is an excellent embedded
database engine used by many projects. Unfortunately the cross
platform compilation required to build Notable on all platforms proved
to be very challenging. Even if it had been feasible it would have
required [cgo](https://golang.org/cmd/cgo/) which complicates the
build and is something this project would like to avoid.

Fortunately Notable is an extremely simple program and adding support
for BoltDB was very easy and solves all of the previously mentioned
challenges. The only challenge it introduces is how to (later) add
proper support for full text searching of notes. Luckily this can
easily be added by leveraging [Bleve](http://www.blevesearch.com/)
which is an excellent full text index for Golang (and doesn't require
cgo).

For platforms other than macos Sqlite3 is still the default backend.
If you would like to try the new backend you can enable it with
`-use.bolt`. Usage of the new BoltDB backend does attempt to migrate
notes from the Sqlite3 format if possible. Migration the other way is
currently not supported (aka the two are not kept in sync).

Additional notes:

* Fix initial subject focus on note creation.
* Improve search performance by simplifying the widget.
* Disable caching to allow clean upgrades, closes [GH-25](https://github.com/jmcfarlane/notable/issues/25).
* Move alerts to the lower right and use upstream class name.
* Initial usage modal triggered by `?`.
* Initial addition of a [BoltDB](https://github.com/boltdb/bolt) backend.
* Usage of `os/user` has been removed in order to enable support for
  macos without the use of cgo.
* The tests have been updated to run twice, once for each backend. The
  intention being to maintain feature parity until a given backend is
  retired.
* Simplfy daemonization logic.
* Add ability to restart via sending a `SIGUSR2` signal.


## v0.0.2 / 2017-03-10

* Include the license in each release.
* Include the version number in the binary.
* Initial tooling to automate the release process.
* Start maintaining a changelog.

## v0.0.1 / 2017-03-09

First release since the project was converted from Python to Golang.
