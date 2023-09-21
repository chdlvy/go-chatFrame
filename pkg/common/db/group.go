package db

import (
	"context"
	"time"
)

type GroupModelInterface interface {
	Create(ctx context.Context, group *GroupModel) (err error)
	Take(ctx context.Context, groupID uint64) (group *GroupModel, err error)
	Find(ctx context.Context, groupIDs []uint64) (groups []*GroupModel, err error)
	// 获取群总数
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
}

// 数据表
type GroupModel struct {
	GroupID       uint64    `gorm:"column:group_id;primary_key;size:64"                 json:"groupID"           binding:"required"`
	GroupName     string    `gorm:"column:name;size:255"                                json:"groupName"`
	FaceURL       string    `gorm:"column:face_url;size:255"                            json:"faceURL"`
	CreatorUserID uint64    `gorm:"column:creator_user_id;size:64"`
	Status        int32     `gorm:"column:status"`
	GroupType     int32     `gorm:"column:group_type"`
	CreateTime    time.Time `gorm:"column:create_time;index:create_time;autoCreateTime"`
}
