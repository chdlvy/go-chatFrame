package push

import (
	"context"
	"github.com/chdlvy/protocol/constant"
	pb "github.com/chdlvy/protocol/push"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

type PushServer struct {
	pb.UnimplementedPushMsgServiceServer
	pusher Pusher
}

func StartPush(PushRpcPort int) {
	listen, _ := net.Listen("tcp", ":"+strconv.Itoa(PushRpcPort))
	grpcServer := grpc.NewServer()
	pb.RegisterPushMsgServiceServer(grpcServer, &PushServer{})
	err := grpcServer.Serve(listen)
	if err != nil {
		log.SetFlags(log.Llongfile)
		log.Fatal("start push rpc failed")
	}
}

func (ps *PushServer) PushMsg(ctx context.Context, req *pb.PushMsgReq) (resp *pb.PushMsgResp, err error) {
	switch req.MsgData.SessionType {
	case constant.GroupChatType:
		err = ps.pusher.Push2Group()
	case constant.PrivateChatType:
		err = ps.pusher.Push2User(req.MsgData)
	}
	return &pb.PushMsgResp{}, err
}
