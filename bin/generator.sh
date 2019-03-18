#!/usr/bin/env bash

# hab https://github.com/go-swagger/go-swagger
# docs https://goswagger.io/generate/spec.html

set -o errexit

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
CONF_PATH=${ROOT}/conf/swagger

__validate() {
    swagger validate ${CONF_PATH}/swagger.yml
}

__swagger1() {
    cd ${ROOT}/api/server/v1
    swagger generate spec -o ${ROOT}/api/server/v1/docs/swagger/swagger.yaml --scan-models
}

__swagger2() {
    cd ${ROOT}/api/server/v2
    swagger generate spec -o ${ROOT}/api/server/v2/docs/swagger/swagger.yaml --scan-models
}

main() {

    case "$1" in
        swagger1)
        __swagger1
        ;;
        swagger2)
        __swagger2
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

  swagger1 - generate an API server

  -h / --help - show this help text and exit 0

EOF
}

main "$@"