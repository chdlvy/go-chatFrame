package hub

import (
	pb "github.com/chdlvy/protocol/hub"
	"google.golang.org/grpc"
	"log"
	"net"
	"server/msggateway"
	"strconv"
)

type HubServer struct {
}

func StartHubServer(port int) {
	listen, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	grpcServer := grpc.NewServer()
	pb.RegisterHubServer(grpcServer, &msggateway.HubServer{})
	err := grpcServer.Serve(listen)
	if err != nil {
		log.SetFlags(log.Llongfile)
		log.Fatal("start hub rpc failed")
	}
}
