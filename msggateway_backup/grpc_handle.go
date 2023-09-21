package msggateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chdlvy/protocol/common"
	"github.com/chdlvy/protocol/msg"
	"github.com/go-playground/validator/v10"
	"server/internal/pkg/rpcclient"
)

type GrpcHandler struct {
	pushRpcClient    *rpcclient.PushRpcClient
	messageRpcClient *rpcclient.MessageRpcClient
	validate         *validator.Validate
}

func NewGrpcHandler(validate *validator.Validate) *GrpcHandler {
	messageRpcClient := rpcclient.NewMessageRpcClient(RpcPort)
	pushRpcClient := rpcclient.NewPushRpcClient(PushRpcPort)
	return &GrpcHandler{
		messageRpcClient: &messageRpcClient,
		pushRpcClient:    &pushRpcClient,
		validate:         validate,
	}
}
func (g *GrpcHandler) SendMsg(ctx context.Context, data *common.MsgData) ([]byte, error) {
	//if err := g.validate.Struct(data); err != nil {
	//	return nil, err
	//}
	msgData := data
	req := &msg.SendMsgReq{MsgData: msgData}
	fmt.Println("grpc_handle")
	//调用msg的rpc方法
	resp, err := g.messageRpcClient.Client.SendMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return res, nil

}
