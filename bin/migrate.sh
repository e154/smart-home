# This file is part of the Smart Home
# Program complex distribution https://github.com/e154/smart-home
# Copyright (C) 2016-2020, Filippov Alex
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

set -o errexit

main() {

    case "${COMMAND}" in
        new)
        __new
        ;;
        up)
        __up
        ;;
        down)
        __down
        ;;
        status)
        __status
        ;;
        gen)
        __gen
        ;;
        help)
        __help
        ;;
        *)
        __help
        exit 1
        ;;
    esac

}

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ENV=development
CONFIG=./conf/dbconfig.yml

MIGRATION_NAME=$2
COMMAND=$1

__new() {
    echo $MIGRATION_NAME

    NOW=$(date +"%Y%m%d_%H%M%S")

    cat  > $BASEDIR/../migrations/"${NOW}_${MIGRATION_NAME}.sql" <<EOF
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


EOF
}

__up() {
    sql-migrate up -config=$CONFIG -env=$ENV
}

__down() {
    sql-migrate down -config=$CONFIG -env=$ENV
}

__status() {
    sql-migrate status -config=$CONFIG -env=$ENV
}

__gen() {

    BD=`which go-bindata`

    if [ -z "$BD" ]; then
      echo "Required tools are missing - check beginning of \"$0\" file for details."
      echo "wait for installing go-bindta"
      go get -u github.com/jteeuwen/go-bindata/...
    fi

    # go get -u github.com/jteeuwen/go-bindata/...
    ${BD} -pkg database -o ${BASEDIR}/../system/migrations/assets/assets.go migrations
}

__help() {
  cat <<EOF
Usage: migrate.sh [options]

OPTIONS:

  new - new migration
  up - up migration
  down - down migration
  status - reset migration
  gen - generate migration sources

  -h / --help - show this help text and exit 0

EOF
}

if [ -z `which sql-migrate` ]; then
  echo "Required tools are missing - check beginning of \"$0\" file for details."
  echo "wait for installing sql-migrate"
  go get github.com/rubenv/sql-migrate/...
fi

main "$@"
