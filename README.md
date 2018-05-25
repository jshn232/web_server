# web_server

## 1.部署proto
将jshn232文件夹复制到$GOPATH/src/目录下
```bash
cp -rf jshn232/ $GOPATH/src/
```
## 2.启动server_go/server.go
```bash
cd server_go
go run server.go
```

## 3.启动grpc客户端（测试用）
```bash
cd grpc_client
go run grpc_client.go
```

## 4.启动客户端（测试用）
```bash
cd client_go
go run client.go
```

此时，服务器端控制台会输出新连接的客户端ip及端口，并通过grpc流推送客户端ip地址列表至grpc客户端