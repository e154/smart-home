#!/usr/bin/env bash

DIR="/opt/smart-home/backup"
NOW=$(date +"%d-%m-%Y-%H-%M-%S")
FILE="${NOW}-dump.sql"
DB_NAME="smarthome_dev"
DB_USER="smarthome"
DB_PASS="smarthome"

mkdir -p ${DIR}

echo "Starting backup to ${FILE}..."
mysqldump -u ${DB_USER} -p"${DB_PASS}" ${DB_NAME} > "${DIR}/${FILE}"