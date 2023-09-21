package db

import (
	"context"
	"time"
)

type GroupMemberModelInterface interface {
	Create(ctx context.Context, groupMembers []*GroupMemberModel) (err error)
	Delete(ctx context.Context, groupID uint64, userIDs []uint64) (err error)
	Take(ctx context.Context, groupID uint64, userID uint64) (groupMember *GroupMemberModel, err error)
	FindGroupMemberIDs(ctx context.Context, groupID uint64) (groupMemberIDs []uint64, err error)
	FindUserJoinedGroupID(ctx context.Context, userID uint64) (groupIDs []uint64, err error)
	TakeGroupMemberNum(ctx context.Context, groupID uint64) (count int64, err error)
}

type GroupMemberModel struct {
	GroupID       uint64    `gorm:"column:group_id;primary_key;size:64"`
	UserID        uint64    `gorm:"column:user_id;primary_key;size:64"`
	NickName      string    `gorm:"column:nickname;size:255"`
	FaceURL       string    `gorm:"column:user_group_face_url;size:255"`
	JoinTime      time.Time `gorm:"column:join_time"`
	InviterUserID uint64    `gorm:"column:inviter_user_id;size:64"`
}
