package rpcclient

import (
	"github.com/chdlvy/protocol/msg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

type Message struct {
	conn   grpc.ClientConnInterface
	Client msg.MsgClient
}
type MessageRpcClient Message

func NewMessage(MsgRpcPort int) *Message {
	conn, err := grpc.Dial("127.0.0.1:"+strconv.Itoa(MsgRpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &Message{
		conn:   conn,
		Client: msg.NewMsgClient(conn),
	}
}
func NewMessageRpcClient(MsgRpcPort int) MessageRpcClient {
	return MessageRpcClient(*NewMessage(MsgRpcPort))
}
