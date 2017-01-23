#!/usr/bin/env bash

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# settings include
settings=${BASEDIR}/settings.sh ; source "$settings" ; if [ $? -ne 0 ] ; then echo "Error - no settings functions $settings" 1>&2 ; exit 1 ; fi

conn="${user}:${password}@tcp(${server})/${base}?charset=utf8&parseTime=true"
export PATH=$PATH:$GOPATH/src/github.com/beego/bee

MIGRATION_NAME=$2
COMMAND=$1

function do_new() {
    echo $MIGRATION_NAME
    bee generate migration $MIGRATION_NAME -fields=""
}

function do_up() {
    bee migrate -driver=${driver} -conn=${conn}
}

function do_down() {
    bee migrate rollback -driver=${driver} -conn=${conn}
}

function do_reset() {
    bee migrate reset -driver=${driver} -conn=${conn}
}

function do_refresh() {
    bee migrate refresh -driver=${driver} -conn=${conn}
}

case "${COMMAND}" in
    new)
    do_new
    ;;
    up)
    do_up
    ;;
    down)
    do_down
    ;;
    reset)
    do_reset
    ;;
    refresh)
    do_refresh
    ;;
    *)
    echo "Usage: $0 new|up|down|reset|refresh" >&2
    exit 3
    ;;
esac