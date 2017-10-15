#!/usr/bin/env bash

set -o errexit

#
# base variables
#
FILENAME=$0
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
EXEC="server"
TMP_DIR="${ROOT}/tmp/${EXEC}"
ARCHIVE="smart-home-${EXEC}.tar.gz"

COMMAND=$*

OS_TYPE="unknown"
OS_ARCH="unknown"

UNAME=`which uname`

if [ -z "$UNAME" ]; then
  message "Required tools are missing - check beginning of \"$0\" file for details."
  exit 1
fi


main() {

    # get os type
    case $COMMAND in
        --init)
            __init
        ;;
        --clean)
            __clean
        ;;
        --build)
            __build
        ;;
        (*)
            __help
        ;;
    esac
}

__init() {

    cd ${ROOT}

    go install $(go list ./... | grep -v "/vendor\|/database")
    go get -u github.com/jteeuwen/go-bindata/...

    cp ${ROOT}/conf/app.sample.conf ${ROOT}/conf/app.conf
    cp ${ROOT}/conf/dev/app.sample.conf ${ROOT}/conf/dev/app.conf
    cp ${ROOT}/conf/dev/db.sample.conf ${ROOT}/conf/dev/db.conf

    __mysql_help
}

__clean() {

    rm -rf ${ROOT}/build
    rm -rf ${ROOT}/tmp
    rm -rf ${ROOT}/vendor/bin
    rm -rf ${ROOT}/vendor/pkg
    rm -rf ${ROOT}/vendor/src
    rm -rf ${TMP_DIR}
}

__build() {

    __check_os

    mkdir -p ${TMP_DIR}

    cd ${ROOT}
    go build -o ${TMP_DIR}/${EXEC}-${OS_TYPE}-${OS_ARCH}

    cp -r ${ROOT}/conf ${TMP_DIR}
    sed 's/dev\/app.conf/prod\/app.conf/' ${ROOT}/conf/app.sample.conf > ${TMP_DIR}/conf/app.sample.conf
    cp ${ROOT}/LICENSE ${TMP_DIR}
    cp ${ROOT}/README.md ${TMP_DIR}
    cp ${ROOT}/contributors.txt ${TMP_DIR}
    cp ${ROOT}/bin/${EXEC} ${TMP_DIR}
    cd ${TMP_DIR}
#    echo "tar: ${ARCHIVE} copy to ${HOME}"
#    tar -zcf ${HOME}/${ARCHIVE} .
}

__mysql_help() {
    cat <<EOF

#
# please install database
#

mysql -u root -p
CREATE DATABASE smarthome;
CREATE USER 'smarthome'@'localhost' IDENTIFIED BY 'smarthome';
GRANT ALL PRIVILEGES ON smarthome . * TO 'smarthome'@'localhost';
FLUSH PRIVILEGES;
use smarthome
source /opt/smart-home/server/dump.sql
EOF
}

__help() {
  cat <<EOF
Usage: project.sh [options]

OPTIONS:

  --init - init project environment in develop mode
  --clean - cleaning of temporary directories
  --build - build documentation

  -h / --help - show this help text and exit 0

EOF
}

__check_os() {

    # get os type
    case `${UNAME} -s` in
        (Linux)
            OS_TYPE="linux"
        ;;
        (Darwin)
            OS_TYPE="darwin"
        ;;
    esac

    # get os arch
    case `${UNAME} -m` in
        (x86_64)
            OS_ARCH="amd64"
        ;;
        (386)
            OS_ARCH="386"
        ;;
        (armv7l)
            OS_ARCH="arm-7"
        ;;
        (armv64l)
            OS_ARCH="arm-64"
        ;;
        (armv6l)
            OS_ARCH="arm-6"
        ;;
        (armv5l)
            OS_ARCH="arm-5"
        ;;
        (mip64)
            OS_ARCH="mip64"
        ;;
        (mip64le)
            OS_ARCH="mip64le"
        ;;
    esac
}

main