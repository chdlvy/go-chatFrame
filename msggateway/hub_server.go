package msggateway

import (
	"fmt"
	"server/constant"
	"server/pkg/common/model"
)

type HubServer struct {
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

func (h *HubServer) OnlinePushMsg(data *model.MsgData) error {
	fmt.Println("hubserver OnlinePushMsg")
	if data.SessionType == constant.SingleChatType {
		//私聊
		return hubServer.LongConnServer.SendMsg(data)
	} else if data.SessionType == constant.GroupChatType {
		//群聊
		return hubServer.LongConnServer.SendGroupMsg(data)
	}
	return nil
}
