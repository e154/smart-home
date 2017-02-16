#!/usr/bin/env bash

set -o errexit

main() {

    case "${COMMAND}" in
        --new)
        __new
        ;;
        --up)
        __up
        ;;
        --down)
        __down
        ;;
        --reset)
        __reset
        ;;
        --refresh)
        __refresh
        ;;
        --help)
        __help
        ;;
        *)
        __help
        exit 1
        ;;
    esac

}

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

settings=${BASEDIR}/settings.sh ; source "$settings" ; if [ $? -ne 0 ] ; then echo "Error - no settings functions $settings" 1>&2 ; exit 1 ; fi
conn="${user}:${password}@tcp(${server})/${base}?charset=utf8&parseTime=true"

MIGRATION_NAME=$2
COMMAND=$1

__new() {
    echo $MIGRATION_NAME
    bee generate migration $MIGRATION_NAME -fields=""
}

__up() {
    bee migrate -driver=${driver} -conn=${conn}
}

__down() {
    bee migrate rollback -driver=${driver} -conn=${conn}
}

__reset() {
    bee migrate reset -driver=${driver} -conn=${conn}
}

__refresh() {
    bee migrate refresh -driver=${driver} -conn=${conn}
}

__help() {
  cat <<EOF
Usage: migrate.sh [options]

OPTIONS:

  --new - new migration
  --up - up migration
  --down - down migration
  --reset - reset migration
  --refresh - refresh migration

  -h / --help - show this help text and exit 0

EOF
}

main "$@"