#!/usr/bin/env bash

#
# sudo apt-get install jq
#

set -o errexit

EMAIL="admin@e154.ru"
PASSWORD="adminadmin"
URL="http://localhost:3000/api/v1"
ACCESS_TOKEN=""
CURRENT_USER=""

JQ=`which jq`

if [ -z "$JQ" ]; then
  message "Required tools are missing - check beginning of \"$0\" file for details."
  message "run command: sudo apt-get install jq"
  exit 1
fi

get_user() {

    __success() {
        CURRENT_USER=${RESULT}
        echo User: ${CURRENT_USER}
    }

    __error() {
        MESSAGE="$( echo ${RESULT} | ${JQ} -r ".message" )"
        echo Message: ${MESSAGE}
    }

    RESULT="$( curl -H "access_token:${ACCESS_TOKEN}" -X GET -s ${URL}/user/{$1} )"
    STATUS="$( echo ${RESULT} | ${JQ} -r ".status" )"

    case ${STATUS} in
    null)
    __success
    ;;
    error)
    __error
    ;;
    *)
    echo "Error: Invalid status"
    echo Status: ${STATUS}
    exit 1
    ;;
    esac
}

get_token() {

    __success() {
        ACCESS_TOKEN="$( echo ${RESULT} | ${JQ} -r ".access_token" )"
        USER_ID="$( echo ${RESULT} | ${JQ} -r ".current_user.id" )"
        echo Your token: ${ACCESS_TOKEN}
        get_user ${USER_ID}
    }

    __error() {
        MESSAGE="$( echo ${RESULT} | ${JQ} -r ".message" )"
        echo Message: ${MESSAGE}
    }

    # https://curl.haxx.se/docs/manpage.html#-u
    RESULT="$( curl -H "Content-Type: application/json; charset=utf-8" -u "${EMAIL}:${PASSWORD}" -s ${URL}/signin )"

    STATUS="$( echo ${RESULT} | ${JQ} -r ".status" )"

    case ${STATUS} in
    null)
    __success
    ;;
    error)
    __error
    ;;
    *)
    echo "Error: Invalid status"
    echo Status: ${STATUS}
    exit 1
    ;;
    esac

}

main() {
    get_token
}

main