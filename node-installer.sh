#!/usr/bin/env bash

#
# sudo apt-get install jq
#
# curl -s http://localhost:1377/smart-home/node-installer.sh | bash /dev/stdin --install
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
#DOWNLOAD_URL="https://github.com/e154/smart-home/releases/download/v0.0.5/smart-home-node.tar.gz"
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

    file="${INSTALL_DIR}/node/conf/node.conf"
    if [ ! -f "$file" ]; then
        log "Create file $file"
        cp ${INSTALL_DIR}/node/conf/node.sample.conf $file
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
