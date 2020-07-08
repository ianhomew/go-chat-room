 # usage
 proto 本身沒有支援 go 的 protobuf 所以需要額外安裝插件
 ```shell script
 go get -u github.com/golang/protobuf/protoc-gen-go
 cp protoc-gen-go  /usr/local/bin/
 ```
 ```shell script
vim ~/.zshrc
  ```
vim 進入後到最底下新增:
```shell script
export GO_PATH=~/go
export PATH=$PATH:$GO_PATH/bin
```
再來
 ```shell script
source ~/.zshrc
  ```




 # command
 
 ```shell script
// 指定單一文件
 protoc ./proto/helloworld.proto --go_out=./src/config  
 ```

 ```shell script
// 全部
 protoc ./proto/*.proto --go_out=./src/config  
 ```



 ```shell script
// js
 protoc ./proto/*.proto --js_out=./src/config  
 ```
