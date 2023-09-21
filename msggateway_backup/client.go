package msggateway

import (
	"context"
	"encoding/json"
	"github.com/chdlvy/protocol/common"
	"golang.org/x/net/websocket"
	"log"
	"strconv"
	"time"
)

const (
	pongMessage = 1 //textFrame
)

type Client struct {
	conn       LongConn
	UserID     uint64
	token      string
	PlatformID int
	grpcHandle *GrpcHandler
}

func NewClient(conn LongConn, userId uint64) *Client {
	return &Client{
		conn:   conn,
		UserID: userId,
	}
}
func (c *Client) ResetClient(conn LongConn) {
	c.conn = conn
}
func (c *Client) readMessage() {
	_ = c.conn.SetReadDeadline(pongWait)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.SetFlags(log.Llongfile)
			log.Println(err)
			return
		}

		log.Print("readMessage = ", string(message))
		data := &common.MsgData{}
		json.Unmarshal(message, data)
		if data.SendID == strconv.Itoa(int(c.UserID)) {
			res, err := c.grpcHandle.SendMsg(context.Background(), data)
			if err != nil {
				log.SetFlags(log.Llongfile)
				log.Fatal(err)
			}
			log.Println("sendMsg res：", string(res))
		}

	}
}
func (c *Client) writeMessage(data []byte) error {
	c.conn.SetWriteDeadline(pongWait)
	return c.conn.WriteMessage(pongMessage, data)
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		err := c.conn.WriteMessage(websocket.TextFrame, []byte("hertbeat"))
		if err != nil {
			log.SetFlags(log.Llongfile)
			log.Fatal(err)
			return
		}
	}
}
