package msggateway

import "time"

const (
	RpcPort     = 5000
	PushRpcPort = 5001
	HubRpcPort  = 5002
)

var Config *configStruct

type configStruct struct {
	LongConnconf struct {
		port                int
		WebsocketMaxConnNum int
		WebsocketMaxMsgLen  int
		WebsocketTimeout    time.Duration
	}
}
