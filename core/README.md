# JudgeSystem
### Client
Dev

### TECH
- MongoDB
- C/C++

### How to build
```sh 
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
```
**Note**: install the lib above, then you can link the lib and include the head files
