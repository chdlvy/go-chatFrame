package db

import (
	"chatFrame/constant"
	"chatFrame/pkg/common/model"
	"gorm.io/gorm"
	"time"
)

type ChatLogGorm struct {
	*MetaDB
}

func NewChatLogGorm(db *gorm.DB) ChatLogModelInterface {
	return &ChatLogGorm{NewMetaDB(db, &ChatLogModel{})}
}
func (c *ChatLogGorm) Create(msg *model.MsgData) error {
	chatLog := &ChatLogModel{}
	switch chatLog.SessionType {
	case constant.GroupChatType:
		chatLog.RecvID = msg.GroupID
	case constant.SingleChatType:
		chatLog.RecvID = msg.RecvID
	}
	if msg.ContentType == constant.TextMsg {
		chatLog.Content = string(msg.Content)
	}
	chatLog.SendTime = time.Unix(0, msg.SendTime*1e6)
	chatLog.CreateTime = time.Unix(0, msg.SendTime*1e6)
	return c.DB.Create(chatLog).Error
}
func (c *ChatLogGorm) DelGroupAllLogWithUser(userID, groupID uint64) error {
	return c.DB.Where("session_type = ? and send_id = ? and recv_id = ?", constant.GroupChatType, userID, groupID).Delete(&ChatLogModel{}).Error
}
