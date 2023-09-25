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

#!/usr/bin/env bash

#
# sudo apt-get install jq
#
# curl -s https://e154.github.io/smart-home/node-installer.sh | bash /dev/stdin --install
#

shopt -s extglob
set -o errtrace
set -o errexit

GIT_USER="e154"
GIT_REPO="smart-home-node"
INSTALL_DIR="/opt/smart-home"
ARCHIVE="node.tar.gz"
DOWNLOAD_URL="$( curl -s https://api.github.com/repos/${GIT_USER}/${GIT_REPO}/releases/latest | jq -r ".assets[].browser_download_url" )"
#DOWNLOAD_URL="https://github.com/e154/smart-home-old/releases/download/v0.0.5/smart-home-node.tar.gz"
COMMAND=$1

JQ=`which jq`

main() {

    case "${COMMAND}" in
        --install)
        __install
        ;;
        --update)
        __update
        ;;
        --node)
        __node
        ;;
        *)
        __help
        exit 1
        ;;
    esac
}

log()  { printf "%b\n" "$*"; }
debug(){ [[ ${node_debug_flag:-0} -eq 0 ]] || printf "%b\n" "Running($#): $*"; }
fail() { log "\nERROR: $*\n" ; exit 1 ; }

__install_initialize() {

    log ""
    log "Install smart home node to: ${INSTALL_DIR}/node"

    if [ -z "$JQ" ]; then
      log "Install jq"
      sudo apt-get install jq
    fi

    log "Trying to install GNU version of tar, might require sudo password"

    sudo mkdir -p ${INSTALL_DIR}
    sudo chown $USER:$USER ${INSTALL_DIR} -R
    mkdir -p ${INSTALL_DIR}/node

    cd ${INSTALL_DIR}/node

    log "Download latest release from:"
    log "URL: ${DOWNLOAD_URL}"
    curl -sSL -o ${ARCHIVE} ${DOWNLOAD_URL}

    log "Unpack archive"
    tar -zxf ${ARCHIVE}
}

__install_default_settings() {

    cd ${INSTALL_DIR}/node

    file="${INSTALL_DIR}/node/conf/config.json"
    if [ ! -f "$file" ]; then
        log "Create file $file"
        cp ${INSTALL_DIR}/node/conf/config.dev.json $file
    fi

}

__install_main() {

#    log "Install node as service"
#    sudo /opt/smart-home/node/node install

    log "node installed"
    exec /opt/smart-home/node/node help

    return 0
}

__install() {
  __install_initialize
  __install_default_settings
  __install_main
}

__update() {

    cd ${INSTALL_DIR}/node

    log "Download latest release from:"
    log "URL: ${DOWNLOAD_URL}"
    curl -sSL -o ${ARCHIVE} ${DOWNLOAD_URL}

    log "Unpack archive"
    tar -zxf ${ARCHIVE}

}

__remove() {

    log "Remove Smart home node"

    log "Stop service"
    sudo ${INSTALL_DIR}/server/node stop

    log "Remove service"
    sudo ${INSTALL_DIR}/server/node remove

    log "Remove node directory"
    rm -rf ${INSTALL_DIR}/node

    return 0
}

__help() {
  cat <<EOF
Usage: node-installer.sh [options]

OPTIONS:

  --install - install node
  --update - update node
  --node - node node

  -h / --help - show this help text and exit 0

EOF
}

main "$@"
