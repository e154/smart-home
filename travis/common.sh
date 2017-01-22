#!/usr/bin/env bash

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONFIGURATOR="${BASEDIR}/../configurator"
NODE="${BASEDIR}/../node"
SERVER="${BASEDIR}/../server"
TMP_DIR="${BASEDIR}/../tmp"












#VERSION_VAR="main.VersionString"
#REV_VAR="main.RevisionString"
#GENERATED_VAR="main.GeneratedString"
#VERSION_VALUE="$(git describe --always --dirty 2> /dev/null)"
#REV_VALUE="$(git rev-parse HEAD 2> /dev/null || echo '???')"
#GENERATED_VALUE="$(date -u +'%Y-%m-%dT%H:%M:%S%z')"

#OS="$(uname | tr '[:upper:]' '[:lower:]')"
#ARCH="$(uname -m | if grep -q x86_64 ; then echo amd64 ; else uname -m ; fi)"