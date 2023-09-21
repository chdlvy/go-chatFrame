package db

import (
	"context"
	"time"
)

type FriendModelInterface interface {
	// 插入多条记录
	Create(ctx context.Context, friends []*FriendModel) (err error)
	// 删除ownerUserID指定的好友
	Delete(ctx context.Context, ownerUserID uint64, friendUserIDs []uint64) (err error)
	// 更新好友信息的非零值
	Update(ctx context.Context, ownerUserID, friends uint64, info map[string]interface{}) (err error)
	// 更新好友备注（也支持零值 ）
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID uint64, remark string) (err error)
	// 获取单个好友信息，如没找到 返回错误
	Take(ctx context.Context, ownerUserID, friendUserID uint64) (friend *FriendModel, err error)
	//获取所有好友的信息
	FindFriends(ctx context.Context, ownerUserID uint64) (friends []*FriendModel, err error)
	//获取所有好友的ID
	FindFriendIDs(ctx context.Context, ownerUserID uint64) (friendIDs []uint64, err error)
}

type FriendModel struct {
	OwnerUserID  uint64    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID uint64    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark       string    `gorm:"column:remark;size:255"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime"`
	AddSource    int32     `gorm:"column:add_source"`
}
