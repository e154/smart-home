#!/usr/bin/env bash

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

${BASEDIR}/make_configurator.sh --init
${BASEDIR}/make_configurator.sh --build-front
${BASEDIR}/make_configurator.sh --build-back
${BASEDIR}/make_node.sh --init
${BASEDIR}/make_node.sh --build
${BASEDIR}/make_tools.sh --init
${BASEDIR}/make_tools.sh --build
${BASEDIR}/make_server.sh --init
${BASEDIR}/make_server.sh --build

