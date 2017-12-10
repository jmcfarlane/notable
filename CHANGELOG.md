## v0.1.3 / ...

* Prevent `./scripts/release.sh` from being called directly

## v0.1.2 / 2017-12-10

* Fix Docker tags on release (was not updating the docker.io tags)
* Improve encryption by using AES-GCM [via cryptopasta](https://github.com/gtank/cryptopasta)
* Migrate to [Dep](https://github.com/golang/dep) for dependency management
* Drop support for Go <= 1.7 (Dep is worth it ;)
* Add integration with https://codecov.io
* Add `make coverage` build target for detail on test coverage
* Fix incorrect db path setup used by tests :/
* Initial tests for secondary nodes
* Minor clean up of version string calculations, add `pid`
* Reduce dead code
* Remove unnecessary `db.migrate()` call from tests :/
* Update dependencies
* Reduce global state (much more to be done here)
* Add `/api/stop` handler and make `/api/restart` require `msg`
* Reduce race conditions (credit: `-race`!)
* Ensure all tests persist to a temporary directory
* Fix defects related to http `handler` error handling (more to do)
* Fix duplicate args when using `-daemon=true -browser=true`

## v0.1.1 / 2017-10-29

* Upgrade to Go 1.9.2
* Initial support for automatic client consumption of changes when
  Notable is being ran in distributed mode. Currently this performs a
  very basic reload of the notes table, but does not take into
  consideration any notes currently open for edit. Reloading works on
  both primary and secondary nodes (via different triggers).
* Upgrade runnable Docker container from Alpine to `scratch`
* Produce more release artifacts via Docker (less required on the host)
* Use `-race` when running tests
* Fix regression in tests (need to clean up from a stalled CLI effort)
* Ensure artifact copy from container to host is rootless ;)
* Update features list to include distributed writes
* Initial (successful) testing with Keybase!


## v0.1.0 / 2017-09-17

* Initial (experimental) support for distributed writes.
    * The idea is you can share your notes directory via a tool like
      Syncthing, and then run a single primary and as many
      `-secondary` instances as you like. All computers have write
      access, and all changes get replicated through the primary back
      to the secondary nodes. Current replication mechanisms (maybe)
      known to work:
        1. [Syncthing](https://syncthing.net/)
    * Testing has only been performed on Linux
    * UI is very much incomplete
* Remove sqlite
* Upgrade dependencies
* Purge the `vendor` folder before a build

## v0.0.10 / 2017-05-15

* Add flag `bolt.timeout` for use with opening BoltDB.
* Move `init` logic into `main` to fix race conditions on startup.
* Open BoltDB before Bleve as the former supports a timeout :)
* Improve error messaging on startup.
* Enable line wrap.
* Enable vendoring via [Glide](https://glide.sh/).
* Ensure all released binaries are built with vendored dependencies.
* Initial work on a command line client [`notable-cli`](https://github.com/jmcfarlane/notable/tree/master/cmd/notable-cli).
* Fix case incorrect import of [logrus](https://github.com/sirupsen/logrus).

## v0.0.9 / 2017-04-26

* Introduce full text index via [Bleve](http://www.blevesearch.com/).
* Add a visual indication if the body of a note has unsaved changes.
* Drop support for <= golang-1.4 as Bleve uses `strings.Compare`.
* Add build branch to `-version` output.
* Fix unhandled error when using the `enter` key with no note selected.
* Add support for FreeBSD.

## v0.0.8 / 2017-03-31

* Allow search to be resumed after manual focus.
* Initial support for (amd64) Microsoft Windows.

## v0.0.7 / 2017-03-29

* Allow search to be resumed after switching tabs.
* Allow forward slashes in subject, tag, and password fields.
* Disable cgo on Linux. Compile from source if you want sqlite.
* Publish a [rkt](https://coreos.com/rkt/) container as a Github release artifact.

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
