package chatlog

import (
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/chdlvy/go-chatFrame/pkg/common/model"
	"log"
)

type ChatLogServer struct {
	chatLogDatabase db.ChatLogDatabase
}

func NewChatLogServer() *ChatLogServer {
	dbconn, err := db.NewGormDB()
	if err != nil {
		log.Println(err)
	}
	chatLogDB := db.NewChatLogGorm(dbconn)
	return &ChatLogServer{chatLogDatabase: db.NewChatLogDatabase(chatLogDB)}
}
func (c *ChatLogServer) CreateChatLog(msg *model.MsgData) error {
	return c.chatLogDatabase.CreateChatLog(msg)
}
func (c *ChatLogServer) DelGroupAllChatLogWithUser(userID, groupID uint64) error {
	return c.chatLogDatabase.DelGroupAllChatLog(userID, groupID)
}
func StartChatLogServer() error {

	//mysql
	gdb, err := db.NewGormDB()
	if err != nil {
		return err
	}
	if err := gdb.AutoMigrate(&db.ChatLogModel{}); err != nil {
		return err
	}
	return nil
}
