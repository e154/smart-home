#!/bin/bash

set -e

echo

UNAME=$(uname -a)
ROOT="./data/vosk"


if [[ "${UNAME}" == *"x86_64"* ]]; then
  ARCH="x86_64"
  echo "amd64 architecture confirmed."
  elif [[ "${UNAME}" == *"aarch64"* ]] || [[ "${UNAME}" == *"arm64"* ]]; then
  ARCH="aarch64"
  echo "aarch64 architecture confirmed."
  elif [[ "${UNAME}" == *"armv7l"* ]]; then
  ARCH="armv7l"
  echo "armv7l (32-bit) WARN: The Coqui and VOSK bindings are broken for this platform at the moment, so please choose Picovoice when the script asks. wire-pod is designed for 64-bit systems."
  STT=""
else
  echo "Your CPU architecture not supported. This script currently supports x86_64, aarch64, and armv7l."
  exit 1
fi

if [[ ${UNAME} == *"Darwin"* ]]; then
  if [[ -f /usr/local/Homebrew/bin/brew ]] || [[ -f /opt/Homebrew/bin/brew ]]; then
    TARGET="darwin"
    echo "macOS detected."
  fi
  elif [[ -f /usr/bin/apt ]]; then
  TARGET="debian"
  echo "Debian-based Linux detected."
  elif [[ -f /usr/bin/pacman ]]; then
  TARGET="arch"
  echo "Arch Linux detected."
  elif [[ -f /usr/bin/dnf ]]; then
  TARGET="fedora"
  echo "Fedora/openSUSE detected."
fi

echo "Getting VOSK assets"
rm -fr ${HOME}/.vosk
mkdir -p ${HOME}/.vosk
cd ${HOME}/.vosk
VOSK_VER="0.3.45"
if [[ ${TARGET} == "darwin" ]]; then
  VOSK_VER="0.3.42"
  VOSK_DIR="vosk-osx-${VOSK_VER}"
  elif [[ ${ARCH} == "x86_64" ]]; then
    VOSK_DIR="vosk-linux-x86_64-${VOSK_VER}"
  elif [[ ${ARCH} == "aarch64" ]]; then
    VOSK_DIR="vosk-linux-aarch64-${VOSK_VER}"
  elif [[ ${ARCH} == "armv7l" ]]; then
    VOSK_DIR="vosk-linux-armv7l-${VOSK_VER}"
fi
VOSK_ARCHIVE="$VOSK_DIR.zip"
wget -q --show-progress --no-check-certificate "https://github.com/alphacep/vosk-api/releases/download/v${VOSK_VER}/${VOSK_ARCHIVE}"
unzip "$VOSK_ARCHIVE"
mv "$VOSK_DIR" libvosk
rm -fr "$VOSK_ARCHIVE"

