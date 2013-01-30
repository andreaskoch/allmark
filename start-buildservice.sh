#!/bin/sh

# Starts a filesystem-wachter which compiles the
# the application everytime a file gets changes
# in the current folder.
# source: https://exyr.org/2011/inotify-run/

# Requirements:
# - inotify-tools (https://github.com/rvoicilas/inotify-tools)

# filename: start-buildservice.sh

FORMAT=$(echo -e "\033[1;33m%w%f\033[0m written")
"$@"
while inotifywait -qre close_write --exclude '(.git)' --format "$FORMAT" .
do
    "$@"
    go install server/server.go
done