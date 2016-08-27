#!/bin/bash

cat ./client.pid | while read LINE
do
    cpid=$LINE
    ps -ef | grep $cpid | grep -v grep
    if [[ $? -ne 0 ]]; then
        echo 'restart'
    fi
    break
done

