#!/usr/bin/env bash

socat -d -d -d -x UNIX-LISTEN:${LOCAL_SOCKET_PATH},fork TCP:${REMOTE}