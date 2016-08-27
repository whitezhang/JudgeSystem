#!/bin/bash
db_backup="/home/judger/worksystem/mongodb_bak"
date=`date +%Y%m%d`

rm -rf ${db_backup}/${date}
mongodump -h 127.0.0.1 -p 27071 --out ${db_backup}/${date}
