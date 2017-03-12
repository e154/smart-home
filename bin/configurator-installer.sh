#!/usr/bin/env bash

#
# sudo apt-get install jq
#
# curl -s http://localhost:1377/smart-home/configurator-installer.sh | bash /dev/stdin --install
# curl -s https://e154.github.io/smart-home/configurator-installer.sh | bash /dev/stdin --install
#

shopt -s extglob
set -o errtrace
set -o errexit

GIT_USER="e154"
GIT_REPO="smart-home-configurator"
INSTALL_DIR="/opt/smart-home"
ARCHIVE="configurator.tar.gz"
DOWNLOAD_URL="$( curl -s https://api.github.com/repos/${GIT_USER}/${GIT_REPO}/releases/latest | jq -r ".assets[].browser_download_url" )"
#DOWNLOAD_URL="https://github.com/e154/smart-home/releases/download/v0.0.5/smart-home-configurator.tar.gz"
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
        *)
        __help
        exit 1
        ;;
    esac
}

log()  { printf "%b\n" "$*"; }
debug(){ [[ ${configurator_debug_flag:-0} -eq 0 ]] || printf "%b\n" "Running($#): $*"; }
fail() { log "\nERROR: $*\n" ; exit 1 ; }

__install_initialize() {

    log ""
    log "Install smart home configurator to: ${INSTALL_DIR}/configurator"

    if [ -z "$JQ" ]; then
      log "Install jq"
      sudo apt-get install jq
    fi

    log "Trying to install GNU version of tar, might require sudo password"

    sudo mkdir -p ${INSTALL_DIR}
    sudo chown $USER:$USER ${INSTALL_DIR} -R
    mkdir -p ${INSTALL_DIR}/configurator

    cd ${INSTALL_DIR}/configurator

    log "Download latest release from:"
    log "URL: ${DOWNLOAD_URL}"
    curl -sSL -o ${ARCHIVE} ${DOWNLOAD_URL}

    log "Unpack archive"
    tar -zxf ${ARCHIVE}
}

__install_default_settings() {

    cd ${INSTALL_DIR}/configurator

    file="${INSTALL_DIR}/configurator/conf/app.conf"
    if [ ! -f "$file" ]; then
        log "Create file $file"
        sed 's/dev\/app.conf/prod\/app.conf/' ${INSTALL_DIR}/configurator/conf/app.sample.conf > $file
    fi

    file="${INSTALL_DIR}/configurator/conf/prod/app.conf"
    if [ ! -f "$file" ]; then
        log "Create file $file"
        cp ${INSTALL_DIR}/configurator/conf/prod/app.sample.conf $file
    fi

}

__install_main() {

#    log "Install configurator as service"
#    sudo /opt/smart-home/configurator/configurator install

    log "configurator installed"
    exec /opt/smart-home/configurator/configurator help

    return 0
}

__install() {
  __install_initialize
  __install_default_settings
  __install_main
}

__update() {

    cd ${INSTALL_DIR}/configurator

    log "Download latest release from:"
    log "URL: ${DOWNLOAD_URL}"
    curl -sSL -o ${ARCHIVE} ${DOWNLOAD_URL}

    log "Unpack archive"
    tar -zxf ${ARCHIVE}

}

__help() {
  cat <<EOF
Usage: configurator-installer.sh [options]

OPTIONS:

  --install - install configurator
  --update - update configurator

  -h / --help - show this help text and exit 0

EOF
}

main "$@"
