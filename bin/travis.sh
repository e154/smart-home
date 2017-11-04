#!/usr/bin/env bash

set -o errexit

#
# base variables
#
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"
EXEC="server"
TMP_DIR="${ROOT}/tmp/${EXEC}"
ARCHIVE="smart-home-${EXEC}.tar.gz"

#
# build version variables
#
PACKAGE="github.com/e154/smart-home"
VERSION_VAR="main.VersionString"
REV_VAR="main.RevisionString"
REV_URL_VAR="main.RevisionURLString"
GENERATED_VAR="main.GeneratedString"
DEVELOPERS_VAR="main.DevelopersString"
BUILD_NUMBER_VAR="main.BuildNumString"
VERSION_VALUE="$(git describe --always --dirty --tags 2>/dev/null)"
REV_VALUE="$(git rev-parse HEAD 2> /dev/null || echo "???")"
REV_URL_VALUE="https://${PACKAGE}/commit/${REV_VALUE}"
GENERATED_VALUE="$(date -u +'%Y-%m-%dT%H:%M:%S%z')"
DEVELOPERS_VALUE="delta54<support@e154.ru>"
BUILD_NUMBER_VALUE="$(echo ${TRAVIS_BUILD_NUMBER})"
GOBUILD_LDFLAGS="\
        -X ${VERSION_VAR}=${VERSION_VALUE} \
        -X ${REV_VAR}=${REV_VALUE} \
        -X ${REV_URL_VAR}=${REV_URL_VALUE} \
        -X ${GENERATED_VAR}=${GENERATED_VALUE} \
        -X ${DEVELOPERS_VAR}=${DEVELOPERS_VALUE} \
        -X ${BUILD_NUMBER_VAR}=${BUILD_NUMBER_VALUE} \
"


main() {

  export DEBIAN_FRONTEND=noninteractive

  : ${INSTALL_MODE:=stable}

  case "$1" in
    --test)
    __test
    ;;
    --init)
    __init
    ;;
    --clean)
    __clean
    ;;
    --migrate)
    __migrate
    ;;
    --docs-deploy)
    __docs_deploy
    ;;
    --build)
    __build
    ;;
    *)
    echo "Error: Invalid argument '$1'" >&2
    exit 1
    ;;
  esac

}

__test() {

   cd ${ROOT}
   goveralls
}

__init() {

    mkdir -p ${TMP_DIR}
    cd ${ROOT}
    gvt rebuild
}

__clean() {

    rm -rf ${ROOT}/vendor/bin
    rm -rf ${ROOT}/vendor/pkg
    rm -rf ${ROOT}/vendor/src
    rm -rf ${TMP_DIR}
    rm -rf ${HOME}/${ARCHIVE}
}

__migrate() {

    sql-migrate up -config=./bin/dbconfig-travis.yml -env="development"

}

__docs_deploy() {

    cd ${ROOT}/doc/themes/default

    npm install
    gulp

    cd ${ROOT}/doc
    hugo

    cd ${ROOT}/doc/public

    git init
    echo -e "Starting to documentation commit.\n"
    git config --global user.email "support@e154.ru"
    git config --global user.name "delta54"

    git remote add upstream "https://$GH_TOKEN@github.com/e154/smart-home.git"
    git fetch upstream
    git reset upstream/gh-pages

    rev=$(git rev-parse --short HEAD)

    git add -A .

    set +o errexit

    git commit -m "rebuild pages at ${rev}"
    git push -q upstream HEAD:gh-pages

    echo -e "Done documentation deploy.\n"

    set -o errexit
}

__build() {

    __docs_deploy

    # build
    cd ${TMP_DIR}

    BRANCH="$(git name-rev --name-only HEAD)"

    if [[ $BRANCH == *"tags/"* ]]; then
      BRANCH="master"
    fi

    echo ""
    echo "build command:"
    echo "xgo --out=${EXEC} --branch=${BRANCH} --targets=linux/*,windows/*,darwin/* --ldflags='${GOBUILD_LDFLAGS}' ${PACKAGE}"
    echo ""

    xgo --out=${EXEC} --branch=${BRANCH} --targets=linux/*,windows/*,darwin/* --ldflags="${GOBUILD_LDFLAGS}" ${PACKAGE}

    # copy configs
    cp -r ${ROOT}/conf ${TMP_DIR}
    #sed 's/dev\/app.conf/prod\/app.conf/' ${ROOT}/conf/app.conf > ${TMP_DIR}/conf/app.conf

    # etc
    cp -r ${ROOT}/examples ${TMP_DIR}
    cp -r ${ROOT}/data ${TMP_DIR}
    cp ${ROOT}/LICENSE ${TMP_DIR}
    cp ${ROOT}/README.md ${TMP_DIR}
    cp ${ROOT}/contributors.txt ${TMP_DIR}

    cp ${ROOT}/bin/server ${TMP_DIR}
    cp ${ROOT}/bin/server-installer.sh ${TMP_DIR}

    # sql dump
    cd ${TMP_DIR}
    mysqldump -u root smarthome > ${TMP_DIR}/dump.sql
    echo "tar: ${ARCHIVE} copy to ${HOME}"

    # create arch
    tar -zcf ${HOME}/${ARCHIVE} .
}

main "$@"