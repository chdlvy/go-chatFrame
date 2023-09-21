package msg

import (
	"context"
	"errors"
	"fmt"
	common "github.com/chdlvy/protocol/common"
	"github.com/chdlvy/protocol/constant"
	pb "github.com/chdlvy/protocol/msg"
	"github.com/chdlvy/protocol/push"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"server/msggateway"
	"strconv"
)

type MsgServer struct {
	pb.UnimplementedMsgServer
	PrivateChatChan chan *common.MsgData
	GroupChatChan   chan *common.MsgData
}

func (ms *MsgServer) PullMessage(ctx context.Context, in *pb.PullMessageReq) (*pb.PullMessageResp, error) {
	return nil, nil
}
func (ms *MsgServer) SearchMessage(ctx context.Context, in *pb.SearchMessageReq) (*pb.SearchMessageResp, error) {
	return nil, nil
}
func (ms *MsgServer) SendMsg(ctx context.Context, req *pb.SendMsgReq) (resp *pb.SendMsgResp, err error) {
	fmt.Println("rpc/msg/server SendMsg")
	resp = &pb.SendMsgResp{}
	if req.MsgData != nil {

		switch req.MsgData.SessionType {
		case constant.PrivateChatType:
			return ms.privateChat(ctx, req)
		case constant.GroupChatType:
			return ms.groupChat(ctx, req)
		default:
			return nil, errors.New("unknown sessionType")
		}
	} else {
		return nil, errors.New("msgData is nil")
	}
}
func (ms *MsgServer) privateChat(ctx context.Context, req *pb.SendMsgReq) (resp *pb.SendMsgResp, err error) {
	fmt.Println("privateChat")
	//ms.PrivateChatChan <- req.MsgData

	conn, err := grpc.Dial("127.0.0.1:"+strconv.Itoa(msggateway.PushRpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := push.NewPushMsgServiceClient(conn)
	//调用push的rpc方法
	_, err = client.PushMsg(ctx, &push.PushMsgReq{MsgData: req.MsgData})
	if err != nil {
		return &pb.SendMsgResp{}, err
	}
	resp = &pb.SendMsgResp{
		SendTime:    req.MsgData.SendTime,
		ClientMsgID: req.MsgData.ClientMsgID,
	}
	return resp, nil
}
func (ms *MsgServer) groupChat(ctx context.Context, req *pb.SendMsgReq) (resp *pb.SendMsgResp, err error) {
	ms.GroupChatChan <- req.MsgData
	resp = &pb.SendMsgResp{
		SendTime:    req.MsgData.SendTime,
		ClientMsgID: req.MsgData.ClientMsgID,
	}
	return resp, nil
}
