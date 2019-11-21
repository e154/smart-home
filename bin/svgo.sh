#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
DIR=${ROOT}/data/icons/*

cd ${ROOT} && svgo ${DIR} --enable=inlineStyles  --config '{ "plugins": [ { "inlineStyles": { "onlyMatchedOnce": false } }] }' --pretty