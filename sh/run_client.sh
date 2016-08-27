#!/bin/bash

echo "$$" > client.pid

while true
do
    dt=`date +%Y%m%d-%H:%M:%S`
    echo $dt 1>>log/log.client
    echo $dt 2>>log/error.client
    ./judgeclient 1>>log/log.client 2>>log/error.client &
    sleep 2
done
