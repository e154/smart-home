#! /bin/sh
### BEGIN INIT INFO
# Provides:          smart-home-server
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Smart home server auto start script
# Description:       http://e154.github.io/smart-home/
### END INIT INFO

# Author: Filippov Alex <filippov.a@e154.ru>

PATH=/sbin:/usr/sbin:/bin:/usr/bin
DESC="Smart Home server auto start script"
NAME=smart-home-server
FOLDER=/opt/smart-home/server
DAEMON=${FOLDER}/server
PIDFILE=/var/run/$NAME.pid
SCRIPTNAME=/etc/init.d/$NAME

# Exit if the package is not installed
[ -x "$DAEMON" ] || exit 0

do_start()
{
  ${DAEMON} > /dev/null 2>&1 &
}

do_stop()
{
  killall ${DAEMON}
  return 0
}

do_status()
{
  return 0
}

do_reload() {
  return 0
}

case "$1" in
  start)
    do_start
    ;;
  stop)
    do_stop
    ;;
  status)
    do_status
    ;;
  restart|force-reload)
    do_stop
    do_start
    ;;
esac

: