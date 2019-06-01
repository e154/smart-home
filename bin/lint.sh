#!/usr/bin/env bash

set -o errexit

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"

__lint() {
    cd ${ROOT} && golangci-lint run --exclude-use-default=false ./...
}

main() {

    case "$1" in
        *)
        __lint
        exit 1
        ;;
    esac

}

main "$@"