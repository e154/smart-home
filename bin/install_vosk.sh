#!/bin/bash

set -e

HOME=$3
ARCH=$2
TARGET=$1

__download() {
  echo "Getting VOSK assets"
  cd ${HOME}
  VOSK_VER="0.3.45"
  if [[ ${TARGET} == "darwin" ]]; then
    VOSK_VER="0.3.42"
    VOSK_DIR="vosk-osx-${VOSK_VER}"
  elif [[ ${ARCH} == "x86" ]]; then
    VOSK_DIR="vosk-linux-x86-${VOSK_VER}"
  elif [[ ${ARCH} == "x86_64" ]]; then
    VOSK_DIR="vosk-linux-x86_64-${VOSK_VER}"
  elif [[ ${ARCH} == "aarch64" ]]; then
    VOSK_DIR="vosk-linux-aarch64-${VOSK_VER}"
  elif [[ ${ARCH} == "armv7l" ]]; then
    VOSK_DIR="vosk-linux-armv7l-${VOSK_VER}"
  elif [[ ${ARCH} == "win64" ]]; then
    VOSK_DIR="vosk-win64-${VOSK_VER}"
  fi
  VOSK_ARCHIVE="$VOSK_DIR.zip"
  wget -q --show-progress --no-check-certificate "https://github.com/alphacep/vosk-api/releases/download/v${VOSK_VER}/${VOSK_ARCHIVE}"
  unzip -o $VOSK_ARCHIVE
  cp -r ${VOSK_DIR}/* ${HOME}
  rm -fr ${VOSK_DIR}
  rm -fr ${VOSK_ARCHIVE}
}

__download
