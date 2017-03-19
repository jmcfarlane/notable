#!/bin/bash -e

cd $(dirname $0)

BACKEND=sqlite3 go test -v -cover -db=/tmp/test.db
BACKEND=boltdb go test -v -cover -db=/tmp/test.db
echo "All tests passed!"

