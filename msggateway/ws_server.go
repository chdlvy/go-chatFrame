package msggateway

import (
	"context"
	"fmt"
	"github.com/chdlvy/go-chatFrame/pkg/common/model"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type LongConnServer interface {
	Run() error
	GetUserConn(userID string) (*Client, bool)
	registerClient(client *Client)
	unregisterClient(client *Client)
	KickUserConn(client *Client)
	SendMsg(data *model.MsgData) (err error)
	SendGroupMsg(data *model.MsgData) (err error)
}

type WsServer struct {
	port             int
	wsMaxConnNum     int
	onlineUserNum    uint64
	handshakeTimeout time.Duration
	clients          *UserMap
	MQ               *MsgQueue
	//后续考虑用grpc解耦
	midServer *MidServer
}

func NewWsServer(wsPort int, wsSocketTimeout time.Duration, wsMaxConnNum int) *WsServer {
	msgqueue, err := NewMsgQueue()
	if err != nil {
		log.Fatal(err)
	}
	if err := msgqueue.InitMQ(); err != nil {
		log.Fatal(err)
	}
	return &WsServer{
		port:             wsPort,
		wsMaxConnNum:     wsMaxConnNum,
		handshakeTimeout: wsSocketTimeout,
		clients:          newUserMap(),
		MQ:               msgqueue,
		//后续考虑用grpc解耦
		midServer: NewMidServer(),
	}
}
func (ws *WsServer) Run() error {
	http.HandleFunc("/", ws.wsHandler)
	return http.ListenAndServe(":"+strconv.Itoa(ws.port), nil)
}

var UserID uint64 = 1

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsLongConn := newGWebSocket(ws.handshakeTimeout)
	err := wsLongConn.GenerateLongConn(w, r)
	fmt.Println("wsLongConn：", wsLongConn)
	wsLongConn.conn.SetCloseHandler(func(code int, text string) error {
		log.Println("frontend conn closed")
		wsLongConn.conn.Close()
		return nil
	})
	if err != nil {
		log.SetFlags(log.Llongfile)
		fmt.Println("Error upgrading connection:", err)
		return
	}
	//初始化客户端
	client := NewClient(wsLongConn, UserID)
	ws.registerClient(client)
	UserID++
	go client.readMessage()
	//go client.heartbeat()
	time.Sleep(1 * time.Second / 2)
	log.Print("websocket start successed,online numbers：", ws.onlineUserNum)
}

// 客户端上线
func (ws *WsServer) registerClient(client *Client) {

	userId := strconv.Itoa(int(client.UserID))
	_, isExist := ws.clients.Get(userId)
	if !isExist {
		ws.clients.Set(userId, client)
		ws.onlineUserNum += 1
	} else {
		//有人顶号就踢掉连接并删除客户端
		ws.KickUserConn(client)
		//重新添加新的客户端
		ws.clients.Set(userId, client)
	}
	//创建一个mq成员
	ws.MQ.CreateMqMember(60000, client)
}

// 客户端下线
func (ws *WsServer) unregisterClient(client *Client) {
	ws.clients.Delete(strconv.Itoa(int(client.UserID)))
	ws.onlineUserNum -= 1
}

// 顶号
func (ws *WsServer) KickUserConn(client *Client) {
	userId := strconv.Itoa(int(client.UserID))
	oldClients, _ := ws.clients.Get(userId)
	//关闭websocket连接并删除client
	oldClients.conn.Close()
	ws.clients.Delete(userId)
}

func (ws *WsServer) GetUserConn(userID string) (*Client, bool) {
	return ws.clients.Get(userID)
}

// 私聊发送消息
func (ws *WsServer) SendMsg(data *model.MsgData) (err error) {
	recvID := strconv.Itoa(int(data.RecvID))
	fmt.Println("ws_server SendMsg")
	client, ok := ws.GetUserConn(recvID)
	if !ok {
		log.Println("push user is not online,userID：", data.RecvID)
	}
	sendData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//消息提醒
	if err := ws.MQ.NotificationPrivateMsg(context.Background(), client, sendData); err != nil {
		return err
	}
	if err := client.writeMessage(int(data.ContentType), sendData); err != nil {
		return err
	}
	//保存聊天记录
	if err := ws.midServer.ChatLogServer.CreateChatLog(data); err != nil {
		return err
	}
	return err
}

// // 在群聊内发送消息
func (ws *WsServer) SendGroupMsg(data *model.MsgData) (err error) {
	//拿到群成员id
	groupMemberIDs, err := ws.midServer.GroupServer.GroupDB.FindGroupMemberIDs(context.Background(), data.GroupID)
	if err != nil {
		return err
	}
	users := make([]*model.MsgData, len(groupMemberIDs))
	for k, v := range groupMemberIDs {
		tmpMsg := data
		tmpMsg.RecvID = v
		users[k] = tmpMsg
	}
	for _, v := range users {
		err := ws.SendMsg(v)
		if err != nil {
			return err
		}
	}
	return nil
}

//
//func (ws *WsServer) PushToAllClients(msg []byte) {
//	clients := make(map[uint64]*Client)
//	ws.clients.GetAllClient(clients)
//	for _, client := range clients {
//		client.writeMessage(msg)
//	}
//}
