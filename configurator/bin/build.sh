#!/bin/bash

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

function do_init() {

    cd ${BASEDIR}/../src
	sudo npm install -g bower
	sudo npm install -g gulp
	npm install
	sudo ln -sr /usr/bin/nodejs  /usr/bin/node
	bower install
	ln -sr src/static_source/bower_components/bootstrap/fonts/ src/static_source/fonts
	gulp
}

function do_clear() {

    cd ${BASEDIR}/../src
	rm -frd node_modules
	rm -frd static_source/bower_components
	rm -frd static_source/css
	rm -frd static_source/js
	rm -frd static_source/tmp
}

case "$1" in
    init)
    do_init
    ;;
    clear)
    do_clear
    ;;
    *)
    echo "Usage: $0 init|clear" >&2
    exit 3
    ;;
esac