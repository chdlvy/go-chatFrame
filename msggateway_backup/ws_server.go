package msggateway

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/chdlvy/protocol/common"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type LongConnServer interface {
	grpc.ClientConnInterface
	Run() error
	GetUserAllCons(userID string) ([]*Client, bool)
	registerClient(client *Client)
	unregisterClient(client *Client)
	KickUserConn(client *Client)
	SendMsg(data *common.MsgData) (err error)
}
type WsServer struct {
	port             int
	wsMaxConnNum     int
	registerChan     chan *Client
	unregisterChan   chan *Client
	onlineUserNum    uint64
	clientPool       sync.Pool
	handshakeTimeout time.Duration
	clients          *UserMap
}

func NewWsServer(config *configStruct) *WsServer {
	return &WsServer{
		port:             config.LongConnconf.port,
		wsMaxConnNum:     config.LongConnconf.WebsocketMaxConnNum,
		handshakeTimeout: config.LongConnconf.WebsocketTimeout,
		registerChan:     make(chan *Client, 1000),
		unregisterChan:   make(chan *Client, 1000),
		clients:          newUserMap(),
		clientPool: sync.Pool{
			New: func() interface{} {
				return new(Client)
			},
		},
	}
}
func (ws *WsServer) Run() error {
	go func() {
		for {
			select {
			case client := <-ws.registerChan:
				ws.registerClient(client)
			case client := <-ws.unregisterChan:
				ws.unregisterClient(client)

			}
		}
	}()
	http.HandleFunc("/", ws.wsHandler)
	return http.ListenAndServe(":"+strconv.Itoa(ws.port), nil)
}
func (ws *WsServer) GetUserAllCons(userID string) ([]*Client, bool) {
	return ws.clients.GetAll(userID)
}

// 客户端上线
func (ws *WsServer) registerClient(client *Client) {
	userId := strconv.Itoa(int(client.UserID))
	_, isExist := ws.clients.Get(userId, client.PlatformID)
	if !isExist {
		ws.clients.Set(userId, client)
		ws.onlineUserNum += 1
	} else {
		//有其他平台的连接存在就踢掉然后重新添加clients
		ws.KickUserConn(client)
		ws.clients.Set(userId, client)
	}
}

// 客户端下线
func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	ws.clients.DeleteAll(strconv.Itoa(int(client.UserID)))
	ws.onlineUserNum -= 1
}

func (ws *WsServer) KickUserConn(client *Client) {
	userId := strconv.Itoa(int(client.UserID))
	oldClients, exist := ws.clients.Get(userId, client.PlatformID)
	if !exist {
		log.SetFlags(log.Llongfile)
		log.Println("not exist client")
		return
	}
	for _, v := range oldClients {
		if v.PlatformID == client.PlatformID {
			v.conn.Close()
			break
		}
	}
	ws.clients.DeleteAll(userId)
}

var UserID uint64 = 1

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsLongConn := newGWebSocket(ws.handshakeTimeout)
	err := wsLongConn.GenerateLongConn(w, r)
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
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(wsLongConn)
	client.UserID = UserID
	client.grpcHandle = NewGrpcHandler(&validator.Validate{})
	UserID++
	ws.registerChan <- client
	go client.readMessage()
	//go client.heartbeat()
	time.Sleep(1 * time.Second / 2)
	log.Print("websocket start successed,online numbers：", ws.onlineUserNum)
}
func (ws *WsServer) SendMsg(data *common.MsgData) (err error) {
	fmt.Println("ws_server SendMsg")
	conns, ok := ws.GetUserAllCons(data.RecvID)
	if !ok {
		log.Println("push user is not online,userID：", data.RecvID)
	}
	for _, client := range conns {
		err = client.writeMessage(data.Content)
	}

	return err
}
