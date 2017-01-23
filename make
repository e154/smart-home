#!/usr/bin/env bash

#set -o errexit

#
# base variables
#
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TMP_DIR="${BASEDIR}/tmp/server"
EXEC=server

#
# build version variables
#
PACKAGE="github.com/e154/smart-home"
VERSION_VAR="main.VersionString"
REV_VAR="main.RevisionString"
REV_URL_VAR="main.RevisionURLString"
GENERATED_VAR="main.GeneratedString"
DEVELOPERS_VAR="main.DevelopersString"
BUILD_NUMBER_VAR="main.BuildNumString"
VERSION_VALUE="$(git describe --always --dirty --tags 2>/dev/null)"
REV_VALUE="$(git rev-parse HEAD 2> /dev/null || echo "???")"
REV_URL_VALUE="https://${PACKAGE}/commit/${REV_VALUE}"
GENERATED_VALUE="$(date -u +'%Y-%m-%dT%H:%M:%S%z')"
DEVELOPERS_VALUE="delta54<support@e154.ru>"
BUILD_NUMBER_VALUE="$(echo ${TRAVIS_BUILD_NUMBER})"
GOBUILD_LDFLAGS="\
        -X ${VERSION_VAR}=${VERSION_VALUE} \
        -X ${REV_VAR}=${REV_VALUE} \
        -X ${REV_URL_VAR}=${REV_URL_VALUE} \
        -X ${GENERATED_VAR}=${GENERATED_VALUE} \
        -X ${DEVELOPERS_VAR}=${DEVELOPERS_VALUE} \
        -X ${BUILD_NUMBER_VAR}=${BUILD_NUMBER_VALUE} \
"


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
        "${BASEDIR}/controllers"
        "${BASEDIR}/models"
        "${BASEDIR}/router"
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
    cd ${BASEDIR}
    gvt rebuild
}

__clean() {

    rm -rf ${BASEDIR}/vendor/bin
    rm -rf ${BASEDIR}/vendor/pkg
    rm -rf ${BASEDIR}/vendor/src
    rm -rf ${TMP_DIR}
}

__build() {

    cd ${BASEDIR}
    go build -ldflags "${GOBUILD_LDFLAGS}" -o ${TMP_DIR}/${EXEC}
    cp -r ${BASEDIR}/conf ${TMP_DIR}
    sed 's/dev\/app.conf/prod\/app.conf/' ${BASEDIR}/conf/app.conf > ${TMP_DIR}/conf/app.conf
}

__help() {
  cat <<EOF
Usage: make [options]

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