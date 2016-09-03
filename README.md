# JudgeSystem
### Version
Dev

### TECH
- Docker
- Golang
- AngularJS
- MongoDB
- C/C++

### How to use
You'd better build this JudgeSystem in a cleaned Centos System, such as Docker.
```sh 
sh -x prepare_env.sh
sh -x build.sh
cd ./output && sh -x ./restart_unserver.sh
```
**Note**: Make sure you have checked the `prepare_env.sh script` and know what it is going to do before running.

### RESTfulAPI(JSON)
#### Signin
```
IP:PORT/slogin?username=xxx&password=xxx
```
- Return 200 if succeed else 400

#### Get the user
```
IP:PORT/suser?uid=xxx
```
- Return 200 if succeed else 400

#### Get the problem
```
IP:PORT/sprobleminfo?pid=xxx
```
- Return 200 if succeed else 400

#### Get the contests
```
IP:PORT/scontests?page=xxx
```
- Return 200 if succeed else 400

#### Get the contest
```
IP:PORT/scontestinfo?cid=xxx
```
- Return 200 if succeed else 400

#### Submit the answer
```
IP:PORT/ssubmit?pid=xxx&code=xxx&lang=xxx
```
- Return 200 if succeed else 400### RESTfulAPI(JSON)
#### Signin
```
IP:PORT/slogin?username=xxx&password=xxx
```
- Return 200 if succeed else 400

#### Get the user
```
IP:PORT/suser?uid=xxx
```
- Return 200 if succeed else 400

#### Get the problem
```
IP:PORT/sprobleminfo?pid=xxx
```
- Return 200 if succeed else 400

#### Get the contests
```
IP:PORT/scontests?page=xxx
```
- Return 200 if succeed else 400

#### Get the contest
```
IP:PORT/scontestinfo?cid=xxx
```
- Return 200 if succeed else 400

#### Submit the answer
```
IP:PORT/ssubmit?pid=xxx&code=xxx&lang=xxx
```
- Return 200 if succeed else 400
