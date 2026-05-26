#!/bin/sh
set -e

cd "$(dirname "$0")/.."

for dir in examples/*/; do
    echo "--- $dir"
    go run "./$dir"
done
