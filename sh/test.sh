#!/bin/bash

source ./send_email.sh
EMAIL_GROUP="whitezhangv5@hotmail.com"
send_email "Exception: offline job" "lbs_client is still running" $EMAIL_GROUP
