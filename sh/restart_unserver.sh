#!/bin/bash

pstree -p | grep "\<UNServer\>"
killall -9 UNServer
sleep 1
./bin/UNServer &
pstree -p | grep "\<UNServer\>"
