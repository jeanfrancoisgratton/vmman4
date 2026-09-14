#!/usr/bin/env sh

set -e

BRANCH=`git rev-parse --abbrev-ref HEAD`
BRANCH=$(echo "$BRANCH" | tr '/' '_')
BINARY=vmman
OUTPUT=/opt/bin
CHECK_PERMS=0
CGO_ENABLED=0

# Parse arguments
while [ "$#" -gt 0 ]; do
    case "$1" in
        -c|--checkperms)
            CHECK_PERMS=1
            ;;
        *)
            OUTPUT="$1"
            ;;
    esac
    shift
done

if [ "$BRANCH" = "master" ] || [ "$BRANCH" = "main" ] || [ "$BRANCH" = "develop" ]; then
    FULLNAME="$BINARY"
else
    FULLNAME="$BINARY-$BRANCH"
fi



echo "Building ${OUTPUT}/${FULLNAME}"
CGO_ENABLED=1 go build -o ${OUTPUT}/${FULLNAME} .
