#!/usr/bin/env bash

#
# sudo apt-get install jq
#

set -o errexit

EMAIL="admin@e154.ru"
PASSWORD="adminadmin"
URL="http://localhost:3000/api/v1"
TOKEN=""
CURRENT_USER=""

get_user() {

    __success() {
        CURRENT_USER=${RESULT}
        echo User: ${CURRENT_USER}
    }

    __error() {
        MESSAGE="$( echo ${RESULT} | jq -r ".message" )"
        echo Message: ${MESSAGE}
    }

    RESULT="$( curl -H "access_token:${TOKEN}" -X GET -s ${URL}/user/{$1} )"
    STATUS="$( echo ${RESULT} | jq -r ".status" )"

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
        TOKEN="$( echo ${RESULT} | jq -r ".token" )"
        USER_ID="$( echo ${RESULT} | jq -r ".current_user.id" )"
        echo Your token: ${TOKEN}
        get_user ${USER_ID}
    }

    __error() {
        MESSAGE="$( echo ${RESULT} | jq -r ".message" )"
        echo Message: ${MESSAGE}
    }

    # https://curl.haxx.se/docs/manpage.html#-u
    RESULT="$( curl -H "Content-Type: application/json; charset=utf-8" -u "${EMAIL}:${PASSWORD}" -s ${URL}/signin )"

    STATUS="$( echo ${RESULT} | jq -r ".status" )"

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