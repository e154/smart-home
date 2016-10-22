#!/usr/bin/env bash

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

function do_init() {

#    go get github.com/beego/bee
    go get github.com/e154/bee
    go get github.com/astaxie/beego
    go get github.com/go-sql-driver/mysql
    go get github.com/astaxie/beego/session/mysql

    return 0
}

function do_test() {

    return 0
}

function do_clear() {

    return 0
}

case "$1" in
    init)
    do_init
    ;;
    test)
    do_test
    ;;
    clear)
    do_clear
    ;;
    *)
    echo "Usage: $0 init|test|clear" >&2
    exit 3
    ;;
esac