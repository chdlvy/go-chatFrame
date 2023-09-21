package msggateway

import (
	"context"
	"fmt"
	pb "github.com/chdlvy/protocol/hub"
	"log"
)

type HubServer struct {
	pb.UnimplementedHubServer
	LongConnServer LongConnServer
}

var hubServer *HubServer

func NewHubServer(server LongConnServer) *HubServer {
	if hubServer == nil {
		hubServer = &HubServer{
			LongConnServer: server,
		}
	}
	return hubServer
}
func (h *HubServer) OnlinePushMsg(ctx context.Context, req *pb.OnlinePushMsgReq) (*pb.OnlinePushMsgResp, error) {
	fmt.Println("hubserver OnlinePushMsg")
	conns, ok := hubServer.LongConnServer.GetUserAllCons(req.MsgData.RecvID)
	if !ok {
		log.Println("push user is not online,userID：", req.MsgData.RecvID)
	}
	fmt.Println(conns)
	for _, client := range conns {
		err := client.writeMessage(req.MsgData.Content)
		if err != nil {
			return nil, err
		}
	}
	return &pb.OnlinePushMsgResp{}, nil
}
