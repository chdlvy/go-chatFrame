package rpc

import (
	pb "github.com/chdlvy/protocol/msg"
	"google.golang.org/grpc"
	"log"
	"net"
	"server/internal/rpc/hub"
	"server/internal/rpc/msg"
	"server/internal/rpc/push"
	"strconv"
)

func StartRpc(rpcPort int, PushRpcPort int, HubRpcPort int) {
	go push.StartPush(PushRpcPort)
	go hub.StartHubServer(HubRpcPort)
	listen, _ := net.Listen("tcp", ":"+strconv.Itoa(rpcPort))
	grpcServer := grpc.NewServer()
	pb.RegisterMsgServer(grpcServer, &msg.MsgServer{})
	err := grpcServer.Serve(listen)
	if err != nil {
		log.SetFlags(log.Llongfile)
		log.Fatal("start rpc failed")
	}
}
