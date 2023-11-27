package msggateway

import (
	"encoding/json"
	"github.com/chdlvy/go-chatFrame/pkg/common/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	conn   LongConn
	UserID uint64
	token  string
}

const (
	pongMessage   = 1 //textFrame
	binaryMessage = 2
)

func NewClient(conn LongConn, userId uint64) *Client {
	return &Client{
		conn:   conn,
		UserID: userId,
	}
}
func (c *Client) readMessage() error {
	_ = c.conn.SetReadDeadline(pongWait)
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.SetFlags(log.Llongfile)
			log.Println("client disconnect：", err)
			return err
		}
		//log.Print("readMessage = ", string(message))

		data := &model.MsgData{}
		//var data *model.MsgData

		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("json.Unmarshal err：", err)
		}

		data.ContentType = int32(messageType)
		if data.IsImage {
			if err := SaveImage(data.Content); err != nil {
				log.Println("save image err", err)
				continue
			}

		}
		if data.SendID == c.UserID {
			//发送消息
			err := hubServer.OnlinePushMsg(data)
			if err != nil {
				log.SetFlags(log.Llongfile)
				log.Println("sendmessage error：", err)
			}
		}
	}
}
func (c *Client) writeMessage(messageType int, data []byte) error {
	c.conn.SetWriteDeadline(pongWait)

	return c.conn.WriteMessage(messageType, data)
}
func SaveImage(data []byte) error {
	var builder strings.Builder
	timestamp := time.Now().Unix()
	builder.WriteString("image_")
	builder.WriteString(strconv.Itoa(int(timestamp)))
	builder.WriteString(".jpg")
	fileName := builder.String()
	img, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer img.Close()
	if _, err := img.Write(data); err != nil {
		return err
	}
	return nil

}
