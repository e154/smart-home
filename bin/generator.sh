#!/usr/bin/env bash

# hab https://github.com/go-swagger/go-swagger
# docs https://goswagger.io/generate/spec.html

set -o errexit

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
CONF_PATH=${ROOT}/conf/swagger

__validate() {
    swagger validate ${CONF_PATH}/${1}/swagger.yml
}

__server() {
    __validate server
    mkdir -p ${ROOT}/api/server_v1/stub
    cd ${ROOT}/api/server_v1/stub
    swagger generate server -f ${CONF_PATH}/server/swagger.yml -A server $@
}

main() {

    case "$1" in
        server)
        __server
        ;;
        *)
        __help
        exit 1
        ;;
    esac

}

__help() {
  cat <<EOF
Usage: generator.sh [options]

OPTIONS:

  server - generate an API server

  -h / --help - show this help text and exit 0

EOF
}

main "$@"