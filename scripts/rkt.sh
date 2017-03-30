#!/bin/bash -e

cd $(dirname $0)/..

source scripts/version.sh

cd target/notable-${TAG}.linux-amd64
acbuild begin
acbuild set-name github.com/jmcfarlane/notable
acbuild copy notable /bin/notable
acbuild set-exec -- /bin/notable -daemon=false -browser=false
acbuild write ../notable-${TAG}.linux-amd64.aci
acbuild end
