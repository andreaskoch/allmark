#!/bin/sh

# Starts a filesystem-wachter which compiles the
# the go application everytime a file gets changes
# in the current folder.

# Source:
# https://exyr.org/2011/inotify-run/

# Requirements:
# inotify-tools (https://github.com/rvoicilas/inotify-tools)

# Filename:
# start-buildservice.sh

# Assemble list of files to watch
autoBuildFileList="autobuild-files.txt"

find `pwd` -name "*.go" > $autoBuildFileList

FORMAT=$(echo -e "\033[1;33m%w%f\033[0m written")
"$@"
while inotifywait -qre close_write --fromfile $autoBuildFileList --format "$FORMAT"
do
    "$@"
    go install
done
