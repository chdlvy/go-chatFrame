package main

import (
	"log"
	"server/internal/rpc"
	"server/msggateway"
	"time"
)

func main() {

	//配置websocket参数
	msggateway.ReSetLongConnConf(8888, 5*time.Second, 100)
	//启动websocket服务
	go rpc.StartRpc(msggateway.RpcPort, msggateway.PushRpcPort, msggateway.HubRpcPort)

	err := msggateway.RunWsAndServer()
	if err != nil {
		log.SetFlags(log.Llongfile)
		log.Println(err)
	}
}
