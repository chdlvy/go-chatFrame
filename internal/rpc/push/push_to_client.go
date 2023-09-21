package push

import (
	"context"
	"fmt"
	"github.com/chdlvy/protocol/common"
	"github.com/chdlvy/protocol/hub"
	"google.golang.org/grpc"
	"server/msggateway"
	"strconv"
)

type Pusher struct {
}

func (p *Pusher) Push2Group() error {
	return nil
}

func (p *Pusher) Push2User(data *common.MsgData) error {
	fmt.Println("push_to_client Push2User")
	conn, _ := grpc.Dial("localhost:"+strconv.Itoa(msggateway.HubRpcPort), grpc.WithInsecure())
	client := hub.NewHubClient(conn)
	//调用hub的rpc方法
	_, err := client.OnlinePushMsg(context.Background(), &hub.OnlinePushMsgReq{MsgData: data})
	if err != nil {
		return err
	}
	return nil
}
