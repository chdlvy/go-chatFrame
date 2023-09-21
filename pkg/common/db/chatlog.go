package db

import (
	"chatFrame/pkg/common/model"
	"time"
)

type ChatLogModel struct {
	SendID         uint64    `gorm:"column:send_id;type:char(64);index:send_id,priority:2"                                                                                         json:"sendID"`
	RecvID         uint64    `gorm:"column:recv_id;type:char(64);index:recv_id,priority:2"                                                                                         json:"recvID"`
	SenderNickname string    `gorm:"column:sender_nick_name;type:varchar(255)"                                                                                                     json:"senderNickname"`
	SenderFaceURL  string    `gorm:"column:sender_face_url;type:varchar(255);"                                                                                                     json:"senderFaceURL"`
	SessionType    int32     `gorm:"column:session_type;index:session_type,priority:2;index:session_type_alone"                                                                    json:"sessionType"`
	MsgFrom        int32     `gorm:"column:msg_from"                                                                                                                               json:"msgFrom"`
	ContentType    int32     `gorm:"column:content_type;index:content_type,priority:2;index:content_type_alone"                                                                    json:"contentType"`
	Content        string    `gorm:"column:content;type:varchar(3000)"                                                                                                             json:"content"`
	Status         int32     `gorm:"column:status"                                                                                                                                 json:"status"`
	SendTime       time.Time `gorm:"column:send_time;index:sendTime;index:content_type,priority:1;index:session_type,priority:1;index:recv_id,priority:1;index:send_id,priority:1" json:"sendTime"`
	CreateTime     time.Time `gorm:"column:create_time"                                                                                                                             json:"ex"`
}
type ChatLogModelInterface interface {
	Create(msg *model.MsgData) error
	DelGroupAllLogWithUser(userID, groupID uint64) error
}
