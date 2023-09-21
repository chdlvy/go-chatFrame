package group

import (
	"server/pkg/common/db"
	"time"
)

type Group struct {
	GroupID       uint64 `gorm:"column:group_id;primary_key;size:64"                 json:"groupID"           binding:"required"`
	GroupName     string `gorm:"column:name;size:255"                                json:"groupName"`
	FaceURL       string `gorm:"column:face_url;size:255"                            json:"faceURL"`
	CreatorUserID uint64 `gorm:"column:creator_user_id;size:64"`
	Status        int32  `gorm:"column:status"`
	GroupType     int32  `gorm:"column:group_type"`
	GroupMember   []*GroupMember
	CreateTime    time.Time `gorm:"column:create_time;index:create_time;autoCreateTime"`
}

type GroupMember struct {
	GroupID       uint64    `gorm:"column:group_id;primary_key;size:64"`
	UserID        uint64    `gorm:"column:user_id;primary_key;size:64"`
	NickName      string    `gorm:"column:nickname;size:255"`
	FaceURL       string    `gorm:"column:user_group_face_url;size:255"`
	JoinTime      time.Time `gorm:"column:join_time"`
	InviterUserID uint64    `gorm:"column:inviter_user_id;size:64"`
}

type GroupInfo struct {
	GroupID       uint64 `protobuf:"bytes,1,opt,name=groupID,proto3" json:"groupID"`
	GroupName     string `protobuf:"bytes,2,opt,name=groupName,proto3" json:"groupName"`
	Introduction  string `protobuf:"bytes,4,opt,name=introduction,proto3" json:"introduction"`
	FaceURL       string `protobuf:"bytes,5,opt,name=faceURL,proto3" json:"faceURL"`
	CreateTime    int64  `protobuf:"varint,7,opt,name=createTime,proto3" json:"createTime"`
	MemberCount   uint32 `protobuf:"varint,8,opt,name=memberCount,proto3" json:"memberCount"`
	Status        int32  `protobuf:"varint,10,opt,name=status,proto3" json:"status"`
	CreatorUserID string `protobuf:"bytes,11,opt,name=creatorUserID,proto3" json:"creatorUserID"`
	GroupType     int32  `protobuf:"varint,12,opt,name=groupType,proto3" json:"groupType"`
}

func Start() error {
	//mysql

	gdb, err := db.NewGormDB()
	if err != nil {
		return err
	}
	if err = gdb.AutoMigrate(&db.GroupModel{}, &db.GroupMemberModel{}); err != nil {
		return err
	}
	return nil
}
