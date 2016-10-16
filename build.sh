#!/usr/bin/env bash

DIRS=(
    '/node/src/serial'
    '/node/src/settings'
)
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

function do_test() {

    for dir in ${DIRS};
    do
        pushd ${BASEDIR}${dir}
        go test -v
        popd
    done

#    go test -v ./node/src/...
}

function do_clear() {

    return 0
}

case "$1" in
    test)
    do_test
    ;;
    clear)
    do_clear
    ;;
    *)
    echo "Usage: $0 test|clear" >&2
    exit 3
    ;;
esac