#!/usr/bin/env bash

set -o errexit

SERVER="/tmp"
TMP_DIR="/tmp"
GOBUILD_LDFLAGS=""
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
common=${BASEDIR}/common.sh ; source "$common" ; if [ $? -ne 0 ] ; then echo "Error - no settings functions $common" 1>&2 ; exit 1 ; fi
GOPATH="${SERVER}/vendor"
TMP_DIR="${TMP_DIR}/server"
EXEC=server

main() {

  export DEBIAN_FRONTEND=noninteractive

  if [[ $# = 0 ]] ; then
    echo 'No arguments provided, installing with'
    echo 'default configuration values.'
  fi

  : ${INSTALL_MODE:=stable}

  case "$1" in
    --test)
    __test
    ;;
    --init)
    __init
    ;;
    --clean)
    __clean
    ;;
    --help)
    __help
    ;;
    --build)
    __build
    ;;
    *)
    echo "Error: Invalid argument '$1'" >&2
    __help
    exit 1
    ;;
  esac

}

__test() {

    DIRS=(
        "${SERVER}/src/api"
    )

    for dir in ${DIRS};
    do
        pushd ${BASEDIR}${dir}
        go test -v
        popd
    done
}

__init() {

    mkdir -p ${TMP_DIR}
    env GOPATH=${GOPATH} go get github.com/FiloSottile/gvt
    cd ${SERVER}
    gvt rebuild
}

__clean() {

    rm -rf ${SERVER}/build
    rm -rf ${SERVER}/vendor/bin
    rm -rf ${SERVER}/vendor/pkg
    rm -rf ${SERVER}/vendor/src
    rm -rf ${TMP_DIR}
}

__build() {

    cd ${SERVER}/src
    env GOPATH=${GOPATH} go build -ldflags "${GOBUILD_LDFLAGS}" -o ${TMP_DIR}/${EXEC}
    cp -r ${SERVER}/src/conf ${TMP_DIR}
    sed 's/dev\/app.conf/prod\/app.conf/' ${SERVER}/src/conf/app.conf > ${TMP_DIR}/conf/app.conf
}

__help() {
  cat <<EOF
Usage: make_server.sh [options]

Bootstrap Debian 8.0 host with mysql installation.

OPTIONS:

  --test - testing package
  --init - initialize the development environment
  --clean - cleaning of temporary directories
  --build - build server

  -h / --help - show this help text and exit 0

EOF
}

main "$@"