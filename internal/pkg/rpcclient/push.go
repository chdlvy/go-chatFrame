package rpcclient

import (
	"github.com/chdlvy/protocol/push"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

type Push struct {
	conn   grpc.ClientConnInterface
	Client push.PushMsgServiceClient
}
type PushRpcClient Push

func NewPush(PushRpcPort int) *Push {
	conn, err := grpc.Dial("127.0.0.1"+strconv.Itoa(PushRpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &Push{
		conn:   conn,
		Client: push.NewPushMsgServiceClient(conn),
	}
}
func NewPushRpcClient(PushRpcPort int) PushRpcClient {
	return PushRpcClient(*NewPush(PushRpcPort))
}
