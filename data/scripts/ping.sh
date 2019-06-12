#!/usr/bin/env bash

ping -c1 $1 > /dev/null 2> /dev/null; [[ $? -eq 0 ]] && echo ok || echo "err"

