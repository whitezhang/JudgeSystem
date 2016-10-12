#!/bin/bash

#cd ./core
#sh -x build.sh

ROOTDIR=`dirname $0 | sed -e "s#^\\([^/]\\)#${PWD}/\\1#"`
GOLANG_VERSION=1.6.1

rm -rf ${ROOTDIR}/go 2>/dev/null
cat go_${GOLANG_VERSION}/go.tar.gz.* | tar xz
cat ${ROOTDIR}/go.tar.gz  | tar zxvf -

export GOROOT=${ROOTDIR}/go
export GOPATH=$ROOTDIR

mkdir -p ${ROOTDIR}/output

rm -rf ${ROOTDIR}/output/* &&

cp -rf ./web/* ./output/ &&

${GOROOT}/bin/go build ./src/UNServer &&

mkdir -p ${ROOTDIR}/output/{conf,bin,log,data,sh} &&

cp ./sh/* ${ROOTDIR}/output/sh &&
cp ./src/UNServer/conf/* output/conf/ &&
#cp ./core/judgeclient output/bin &&
cp UNServer output/bin || exit 1
