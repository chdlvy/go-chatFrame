package db

import "server/pkg/common/model"

type ChatLogDatabase interface {
	CreateChatLog(msg *model.MsgData) error
	DelGroupAllChatLog(userID, groupID uint64) error
}

func NewChatLogDatabase(chatLogModelInterface ChatLogModelInterface) ChatLogDatabase {
	return &chatLogDatabase{chatLogModel: chatLogModelInterface}
}

type chatLogDatabase struct {
	chatLogModel ChatLogModelInterface
}

func (c *chatLogDatabase) CreateChatLog(msg *model.MsgData) error {
	return c.chatLogModel.Create(msg)
}
func (c *chatLogDatabase) DelGroupAllChatLog(userID, groupID uint64) error {
	return c.chatLogModel.DelGroupAllLogWithUser(userID, groupID)
}
