#
# This file is part of the Smart Home
# Program complex distribution https://github.com/e154/smart-home
# Copyright (C) 2016-2023, Filippov Alex
#
# This library is free software: you can redistribute it and/or
# modify it under the terms of the GNU Lesser General Public
# License as published by the Free Software Foundation; either
# version 3 of the License, or (at your option) any later version.
#
# This library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
# Library General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public
# License along with this library.  If not, see
# <https://www.gnu.org/licenses/>.
#

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

    cp ${ROOT}/conf/config.dev.json ${ROOT}/conf/config.json
    cp ${ROOT}/conf/dbconfig.dev.yml ${ROOT}/conf/dbconfig.yml

    __pg_help
}

__clean() {

    rm -rf ${ROOT}/build
    rm -rf ${ROOT}/tmp
    rm -rf ${ROOT}/vendor
    rm -rf ${TMP_DIR}
}

__build() {

    __check_os

    mkdir -p ${TMP_DIR}/conf
    mkdir -p ${TMP_DIR}/data

    cd ${ROOT}
    go build -o ${TMP_DIR}/${EXEC}-${OS_TYPE}-${OS_ARCH}


    cp ${ROOT}/conf/config.dev.json ${TMP_DIR}/conf
    cp ${ROOT}/conf/dbconfig.dev.yml ${TMP_DIR}/conf
    cp ${ROOT}/LICENSE ${TMP_DIR}
    cp ${ROOT}/README* ${TMP_DIR}
    cp ${ROOT}/contributors.txt ${TMP_DIR}
    cp ${ROOT}/bin/server ${TMP_DIR}
    cp -r ${ROOT}/data/icons ${TMP_DIR}/data
    cp -r ${ROOT}/data/scripts ${TMP_DIR}/data
    cp -r ${ROOT}/assets ${TMP_DIR}
    cp -r ${ROOT}/snapshots ${TMP_DIR}
    cd ${TMP_DIR}
#    echo "tar: ${ARCHIVE} copy to ${HOME}"
#    tar -zcf ${HOME}/${ARCHIVE} .
}

__pg_help() {
    cat <<EOF

#
# please install postgres database
#

su - postgres

createuser smart_home;
create database smart_home;

psql ( enter the password for postgressql)

alter user smart_home with encrypted password 'smart_home';
grant all privileges on database smart_home to smart_home;

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
            OS_TYPE="darwin-10.6"
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
