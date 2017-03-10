#!/bin/bash -e

cd $(dirname $0)

go test -v -cover -db=/tmp/test.db
echo "All tests passed!"

