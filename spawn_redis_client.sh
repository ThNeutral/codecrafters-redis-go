#!/bin/sh
set -e
tmpFile=$(mktemp)
go build -o "$tmpFile" client/*.go
exec "$tmpFile" "$@"
