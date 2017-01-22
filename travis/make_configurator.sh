#!/usr/bin/env bash

set -o errexit

CONFIGURATOR="/tmp"
TMP_DIR="/tmp"
GOBUILD_LDFLAGS=""
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
common=${BASEDIR}/common.sh ; source "$common" ; if [ $? -ne 0 ] ; then echo "Error - no settings functions $common" 1>&2 ; exit 1 ; fi
GOPATH="${CONFIGURATOR}/vendor"
TMP_DIR="${TMP_DIR}/configurator"
EXEC=configurator

main() {

  export DEBIAN_FRONTEND=noninteractive

  if [[ $# = 0 ]] ; then
    echo 'No arguments provided, installing with'
    echo 'default configuration values.'
  fi

  : ${INSTALL_MODE:=stable}

  case "$1" in
    test)
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
    --build-front)
    __build_front
    ;;
    --build-back)
    __build_back
    ;;
    *)
    echo "Error: Invalid argument '$1'" >&2
    __help
    exit 1
    ;;
  esac

}

__test() {

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
    cd ${CONFIGURATOR}
    gvt rebuild
    cd ${CONFIGURATOR}/static_source/private
    bower install
    cd ${CONFIGURATOR}/static_source/public
    bower install
    cd ${CONFIGURATOR}/static_source
    npm install
}

__clean() {

    rm -rf ${CONFIGURATOR}/build
    rm -rf ${CONFIGURATOR}/src/tmp
    rm -rf ${CONFIGURATOR}/static_source/node_modules
    rm -rf ${CONFIGURATOR}/static_source/private/bower_components
    rm -rf ${CONFIGURATOR}/static_source/private/tmp
    rm -rf ${CONFIGURATOR}/static_source/public/bower_components
    rm -rf ${CONFIGURATOR}/static_source/public/tmp
    rm -rf ${CONFIGURATOR}/vendor/bin
    rm -rf ${CONFIGURATOR}/vendor/pkg
    rm -rf ${CONFIGURATOR}/vendor/src
    rm -rf ${TMP_DIR}
}

__build_front() {

    cd ${CONFIGURATOR}/static_source
    gulp pack
    cp -r ${CONFIGURATOR}/build ${TMP_DIR}
}

__build_back() {

    cd ${CONFIGURATOR}/src
    env GOPATH=${GOPATH} go build -ldflags "${GOBUILD_LDFLAGS}" -o ${TMP_DIR}/${EXEC}
    cp -r ${CONFIGURATOR}/src/views ${TMP_DIR}
    cp -r ${CONFIGURATOR}/src/conf ${TMP_DIR}
    sed 's/dev\/app.conf/prod\/app.conf/' ${CONFIGURATOR}/src/conf/app.conf > ${TMP_DIR}/conf/app.conf
}

__help() {
  cat <<EOF
Usage: make_configurator.sh [options]

Bootstrap Debian 8.0 host with mysql installation.

OPTIONS:

  --init - initialize the development environment
  --clean - cleaning of temporary directories
  --build-front - build frontend
  --build-back - build backend

  -h / --help - show this help text and exit 0

EOF
}

main "$@"