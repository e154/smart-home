# This file is part of the Smart Home
# Program complex distribution https://github.com/e154/smart-home
# Copyright (C) 2016-2020, Filippov Alex
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

#!/usr/bin/env bash

# hab https://github.com/go-swagger/go-swagger
# docs https://goswagger.io/generate/spec.html

set -o errexit

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
CONF_PATH=${ROOT}/conf/swagger
#SWAGGER=swagger_darwin_amd64.dms.v0.19
SWAGGER=swagger_darwin_amd64.dms.v0.20.1

__validate() {
    ${SWAGGER} validate ${CONF_PATH}/swagger.yml
}

__swagger1() {
    cd ${ROOT}/api/server/v1
    ${SWAGGER} generate spec -o ${ROOT}/api/server/v1/docs/swagger/swagger.yaml --scan-models
}

__swagger2() {
    cd ${ROOT}/api/server/v2
    ${SWAGGER} generate spec -o ${ROOT}/api/server/v2/docs/swagger/swagger.yaml --scan-models
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

  swagger1 - generate an API server v1
  swagger2 - generate an API server v2

  -h / --help - show this help text and exit 0

EOF
}

main "$@"
