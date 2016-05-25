ROOTDIR=`pwd`

export GOPATH=$ROOTDIR

#rm -rf ${ROOTDIR}/output &&

go build ./src/UNServer &&

mkdir -p ${ROOTDIR}/output/{conf,bin,log,data,sh} &&

cp ./src/UNServer/conf/* output/conf/ &&
cp UNServer output/bin || exit 1
