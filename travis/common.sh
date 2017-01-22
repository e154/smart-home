#!/usr/bin/env bash

#
# base variables
#
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONFIGURATOR="${BASEDIR}/../configurator"
NODE="${BASEDIR}/../node"
SERVER="${BASEDIR}/../server"
TMP_DIR="${BASEDIR}/../tmp"

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