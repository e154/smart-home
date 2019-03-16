#!/usr/bin/env bash

# hab https://github.com/go-swagger/go-swagger
# docs https://goswagger.io/generate/spec.html

set -o errexit

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
CONF_PATH=${ROOT}/conf/swagger

__validate() {
    swagger validate ${CONF_PATH}/swagger.yml
}

__swagger() {
    cd ${ROOT}/api/server/v1
    swagger generate spec -o ${ROOT}/api/server/v1/docs/swagger/swagger.yaml --scan-models
}

main() {

    case "$1" in
        swagger)
        __swagger
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

  swagger - generate an API server

  -h / --help - show this help text and exit 0

EOF
}

main "$@"