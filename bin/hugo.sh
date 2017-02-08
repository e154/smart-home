#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"

_hugoserver_ru() {
  cd ${ROOT}/doc
  hugo server --buildDrafts --verbose --source="${ROOT}/doc" --config="${ROOT}/doc/config.yaml" --port=1377 &
}

_hugoserver_en() {
  cd ${ROOT}/doc
  hugo server --buildDrafts --verbose --source="${ROOT}/doc" --config="${ROOT}/doc/config_en.toml" --port=1399 &
}

_hugoserver() {
  cd ${ROOT}/doc
  hugo server --buildDrafts --verbose --source="${ROOT}/doc" --config="${ROOT}/doc/config.toml" --port=1399
}

main() {
    _hugoserver_ru
    _hugoserver_en
#    _hugoserver
}

main "$@"