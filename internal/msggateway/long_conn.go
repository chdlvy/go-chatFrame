package msggateway

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type LongConn interface {
	Close() error
	WriteMessage(messageType int, message []byte) error
	ReadMessage() (int, []byte, error)
	SetReadDeadline(timeout time.Duration) error
	SetWriteDeadline(timeout time.Duration) error
	Dial(urlStr string, requestHeader http.Header) (*http.Response, error)
	GenerateLongConn(w http.ResponseWriter, r *http.Request) error
}

// 生成长连接
type GWebSocket struct {
	conn             *websocket.Conn
	handshakeTimeout time.Duration
}

func newGWebSocket(handshakeTimeout time.Duration) *GWebSocket {
	return &GWebSocket{handshakeTimeout: handshakeTimeout}
}
func (g *GWebSocket) GenerateLongConn(w http.ResponseWriter, r *http.Request) error {
	upgrade := &websocket.Upgrader{
		HandshakeTimeout: g.handshakeTimeout,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade.Upgrade err：", err)
		return err
	}
	g.conn = conn
	return nil
}
func (g *GWebSocket) Close() error {
	return g.conn.Close()
}
func (g *GWebSocket) SetReadDeadline(timeout time.Duration) error {
	return g.conn.SetReadDeadline(time.Now().Add(timeout))
}
func (g *GWebSocket) SetWriteDeadline(timeout time.Duration) error {
	return g.conn.SetWriteDeadline(time.Now().Add(timeout))
}
func (g *GWebSocket) ReadMessage() (int, []byte, error) {
	return g.conn.ReadMessage()
}
func (g *GWebSocket) WriteMessage(messageType int, message []byte) error {
	return g.conn.WriteMessage(messageType, message)
}
func (g *GWebSocket) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err == nil {
		g.conn = conn
	}
	return httpResp, nil
}
