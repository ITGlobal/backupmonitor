#!/usr/bin/env bash
set -e

for gomod in `find -type f -name go.mod`; do
    dir=$(dirname $gomod)
    echo "[$dir]"
    pushd $dir
    go get
    popd
done
