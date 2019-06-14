#!/usr/bin/env bash

#
# sudo apt-get install jq
#
# curl -s https://e154.github.io/smart-home/server-installer.sh | bash /dev/stdin --install
#

shopt -s extglob
set -o errtrace
set -o errexit

GIT_USER="e154"
GIT_REPO="smart-home"
INSTALL_DIR="/opt/smart-home"
ARCHIVE="server.tar.gz"
DOWNLOAD_URL="$( curl -s https://api.github.com/repos/${GIT_USER}/${GIT_REPO}/releases/latest | jq -r ".assets[].browser_download_url" )"
#DOWNLOAD_URL="https://github.com/e154/smart-home-old/releases/download/v0.0.4/smart-home-server.tar.gz"
COMMAND=$1

OS_TYPE="unknown"
OS_ARCH="unknown"

JQ=`which jq`

main() {

    case "${COMMAND}" in
        --install)
        __install
        ;;
        --update)
        __update
        ;;
        --remove)
        __remove
        ;;
        *)
        __help
        exit 1
        ;;
    esac
}

log()  { printf "%b\n" "$*"; }
debug(){ [[ ${server_debug_flag:-0} -eq 0 ]] || printf "%b\n" "Running($#): $*"; }
fail() { log "\nERROR: $*\n" ; exit 1 ; }

__install_initialize() {

    log ""
    log "Install smart home server to: ${INSTALL_DIR}/server"

    if [ -z "$JQ" ]; then
      log "Install jq"
      sudo apt-get install jq
    fi

    log "Trying to install GNU version of tar, might require sudo password"

    sudo mkdir -p ${INSTALL_DIR}
    sudo chown $USER:$USER ${INSTALL_DIR} -R
    mkdir -p ${INSTALL_DIR}/server

    cd ${INSTALL_DIR}/server

    log "Download latest release from:"
    log "URL: ${DOWNLOAD_URL}"
    curl -sSL -o ${ARCHIVE} ${DOWNLOAD_URL}

    log "Unpack archive"
    tar -zxf ${ARCHIVE}
}

__install_default_settings() {

    cd ${INSTALL_DIR}/server

    file="${INSTALL_DIR}/server/conf/config.json"
    if [ ! -f "$file" ]; then
        log "Create file $file"
        cp ${INSTALL_DIR}/server/conf/config.dev.json $file
    fi

    file="${INSTALL_DIR}/server/conf/dbconfig.yml"
    if [ ! -f "$file" ]; then
        log "Create file $file"
        cp ${INSTALL_DIR}/server/conf/dbconfig.dev.yml $file
    fi

    if [ ! -d "${INSTALL_DIR}/data" ]; then
        log "Move data directory ../"
        mv data ../
    fi

    if [ ! -d "${INSTALL_DIR}/assets" ]; then
        log "Move assets directory ../"
        mv assets ../
    fi

    if [ ! -d "${INSTALL_DIR}/snapshots" ]; then
        log "Move snapshots directory ../"
        mv snapshots ../
    fi

}

__install_main() {

#    log "Install server as service"
#    sudo /opt/smart-home/server/server install

    log "Server installed"
    exec /opt/smart-home/server/server help

    return 0
}

__install() {
  __install_initialize
  __install_default_settings
  __install_main
}

__update() {

    cd ${INSTALL_DIR}/server

    log "Download latest release from:"
    log "URL: ${DOWNLOAD_URL}"
    curl -sSL -o ${ARCHIVE} ${DOWNLOAD_URL}

    log "Unpack archive"
    tar -zxf ${ARCHIVE}

    log "Move data directory ../"
    mv assets ../
    mv data/icons ../icons
    mv data/scripts ../scripts
    mv snapshots ../
}

__remove() {

    log "Remove Smart home server"

    log "Stop service"
    sudo ${INSTALL_DIR}/server/server stop

    log "Remove service"
    sudo ${INSTALL_DIR}/server/server remove

    log "Remove server directory"
    rm -rf ${INSTALL_DIR}/server

    return 0
}

__help() {
  cat <<EOF
Usage: server-installer.sh [options]

OPTIONS:

  --install - install server
  --update - update server
  --remove - remove server

  -h / --help - show this help text and exit 0

EOF
}

__check_os() {

    # get os type
    case `${UNAME} -s` in
        (Linux)
            OS_TYPE="linux"
        ;;
        (Darwin)
            OS_TYPE="darwin"
            INSTALL_DIR="${HOME}/smart-home"
        ;;
    esac
}

__check_os

main "$@"
