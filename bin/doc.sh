#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
HUGO=`which hugo`

main() {
  case "$1" in
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

__build() {
    cd ${ROOT}/doc
    hugo
}

__clean() {
    rm -rf ${ROOT}/doc/public
}

__dev() {
  cd ${ROOT}/doc
  hugo server --buildDrafts --verbose --source="${ROOT}/doc" --config="${ROOT}/doc/config.yaml" --port=1377 --disableFastRender
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