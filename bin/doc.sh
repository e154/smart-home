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

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
HUGO=`which hugo`

main() {
  case "$1" in
    --init)
    __init
    ;;
    --clean)
    __clean
    ;;
    --help)
    __help
    ;;
    --dev)
    __dev
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

__init() {
  mkdir -p ${ROOT}/doc
  cd ${ROOT}/doc
  hugo new site
}

__build() {
    cd ${ROOT}/doc
    npm install postcss-cli
    hugo
}

__clean() {
    rm -rf ${ROOT}/doc/public
}

__dev() {
  cd ${ROOT}/doc
  hugo server --buildDrafts --verbose --source="${ROOT}/doc" --config="${ROOT}/doc/config.toml" --port=1377 --disableFastRender
}

__help() {
  cat <<EOF
Usage: doc.sh [options]

OPTIONS:

  --dev - run in develop mode
  --clean - cleaning of temporary directories
  --build - build documentation

  -h / --help - show this help text and exit 0

EOF
}

if [ -z "$HUGO" ]; then
  go get -u -v github.com/gohugoio/hugo
fi

main "$@"
