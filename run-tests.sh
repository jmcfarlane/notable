#!/bin/bash -e

cd $(dirname $0)

go test -v -cover -race -db=/tmp/test.db
echo "All tests passed!"

