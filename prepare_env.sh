yum update
yum install -y boost boost-devel boost-doc
yum install -y gcc-c++ gcc
yum install -y scons
yum install -y git
yum install glibc glibc-devel glibc-static

mkdir -p /home/workspace/
mkdir -p /home/workspace/runtime
mkdir -p /home/workspace/data
cd /home/workspace
git clone -b legacy https://github.com/mongodb/mongo-cxx-driver.git
cd ./mongo-cxx-driver

mkdir -p /usr/local/mongo
scons --prefix=/usr/local/mongo install

# wget https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-3.2.7.tgz
# g++ tutorial.cpp -pthread -lmongoclient -lboost_thread-mt -lboost_system -lboost_regex -o tutorial -I/usr/local/mongo/include -L/usr/local/mongo/lib

	