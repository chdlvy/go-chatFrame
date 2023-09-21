package db

import (
	"context"
	"time"
)

type FriendRequestModelInterface interface {
	// 插入多条申请记录
	Create(ctx context.Context, friendRequests []*FriendRequestModel) (err error)
	//更新记录
	Update(ctx context.Context, friendRequest *FriendRequestModel) (err error)
	// 删除记录
	Delete(ctx context.Context, fromUserID, toUserID uint64) (err error)
	//查找申请记录
	Take(ctx context.Context, fromUserID, toUserID uint64) (friendRequest *FriendRequestModel, err error)
	// 获取toUserID收到的好友申请列表
	FindToUserID(ctx context.Context, toUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error)
	// 获取fromUserID发出去的好友申请列表
	FindFromUserID(ctx context.Context, fromUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error)
}

type FriendRequestModel struct {
	FromUserID    uint64    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      uint64    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time; autoCreateTime"`
	HandlerUserID uint64    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
}
