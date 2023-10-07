package model

type MsgData struct {
	SendID         uint64 `json:"sendID,omitempty"`
	RecvID         uint64 `json:"recvID,omitempty"`
	GroupID        uint64 `json:"groupID,omitempty"`
	SenderNickname string `json:"senderNickname,omitempty"`
	SenderFaceURL  string `json:"senderFaceURL,omitempty"`
	ContentType    int32  `json:"contentType,omitempty"`
	Content        []byte `json:"content,omitempty"`
	SendTime       int64  `json:"sendTime,omitempty"`
	IsRead         bool   `json:"isRead,omitempty"`
	SessionType    int    `json:"sessionType"`
	IsImage        bool   `json:"isImage"`
	Seq            int64  `json:"seq"`
}

type User struct {
	UserID     uint64
	NickName   string
	FaceURL    string
	CreateTime int64
}
