#!/usr/bin/env bash

set -o errexit

TOOLS="/tmp"
TMP_DIR="/tmp"
GOBUILD_LDFLAGS=""
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
common=${BASEDIR}/common.sh ; source "$common" ; if [ $? -ne 0 ] ; then echo "Error - no settings functions $common" 1>&2 ; exit 1 ; fi
GOPATH="${TOOLS}/vendor"
TMP_DIR="${TMP_DIR}/tools"

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
        "${TOOLS}/src/controllers"
        "${TOOLS}/src/models"
        "${TOOLS}/src/router"
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
}

__clean() {

    rm -rf ${TMP_DIR}
}

__build() {

    cd ${TOOLS}/test_node_modbus/src
    env GOPATH=${GOPATH} go build -o ${TMP_DIR}/modbus_tester
    cp ${TOOLS}/test_node_modbus/src/node_modbus.conf ${TMP_DIR}

    cd ${TOOLS}/test_node_serial/src
    env GOPATH=${GOPATH} go build -o ${TMP_DIR}/serial_tester
    cp ${TOOLS}/test_node_serial/src/node_serial.conf ${TMP_DIR}

    cd ${TOOLS}/performance/src
    env GOPATH=${GOPATH} go build -o ${TMP_DIR}/performance
    cp ${TOOLS}/performance/src/performance.conf ${TMP_DIR}
}

__help() {
  cat <<EOF
Usage: make_tools.sh [options]

Bootstrap Debian 8.0 host with mysql installation.

OPTIONS:

  --test - testing package
  --init - initialize the development environment
  --clean - cleaning of temporary directories
  --build - build backend

  -h / --help - show this help text and exit 0

EOF
}

main "$@"