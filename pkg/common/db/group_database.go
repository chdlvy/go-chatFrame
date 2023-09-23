package db

import (
	"context"
	"github.com/chdlvy/go-chatFrame/pkg/common/db/cache"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type GroupDatabase interface {
	//group method
	CreateGroup(ctx context.Context, group *GroupModel, groupMembers []*GroupMemberModel) error
	TakeGroup(ctx context.Context, groupID uint64) (group *GroupModel, err error)
	FindGroup(ctx context.Context, groupIDs []uint64) (groups []*GroupModel, err error)
	JoinGroup(ctx context.Context, groupID, inviterID uint64, users []*UserModel) error
	// 获取群总数
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)

	//GroupMember method
	CreateGroupMember(ctx context.Context, groupMembers []*GroupMemberModel) (err error)
	DeleteGroupMember(ctx context.Context, groupID uint64, userIDs []uint64) error
	TakeGroupMember(ctx context.Context, groupID uint64, userID uint64) (groupMember *GroupMemberModel, err error)
	FindGroupMemberIDs(ctx context.Context, groupID uint64) (groupMemberIDs []uint64, err error)
	FindUserJoinedGroupID(ctx context.Context, userID uint64) (groupIDs []uint64, err error)
	TakeGroupMemberNum(ctx context.Context, groupID uint64) (count int64, err error)
	//other method

}

func InitGroupDatabase(db *gorm.DB) GroupDatabase {
	groupDB := NewGroupDB(db)
	groupMemberDB := NewGroupMemberDB(db)
	rdb, err := cache.NewRedis()
	if err != nil {
		log.Fatal(err)
	}
	groupRDB := cache.NewGroupCache(rdb)
	return &groupDatabase{
		groupDB:     groupDB,
		groupMember: groupMemberDB,
		groupRDB:    groupRDB,
		tx:          db,
	}
}

type groupDatabase struct {
	groupDB     GroupModelInterface
	groupMember GroupMemberModelInterface
	groupRDB    cache.GroupCache
	tx          *gorm.DB
}

func (g *groupDatabase) CreateGroup(ctx context.Context, group *GroupModel, groupMembers []*GroupMemberModel) error {
	//先尝试从缓存中获取，获取不到再拿数据库
	//暂时先实现操作数据库
	if err := g.tx.Transaction(func(tx *gorm.DB) error {
		if err := g.groupDB.Create(ctx, group); err != nil {
			return err
		}
		if len(groupMembers) > 0 {
			if err := g.groupMember.Create(ctx, groupMembers); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil

}

func (g *groupDatabase) TakeGroup(ctx context.Context, groupID uint64) (group *GroupModel, err error) {
	//mysql
	return g.groupDB.Take(ctx, groupID)

}

func (g *groupDatabase) FindGroup(ctx context.Context, groupIDs []uint64) (groups []*GroupModel, err error) {
	//mysql
	return g.groupDB.Find(ctx, groupIDs)
}

func (g *groupDatabase) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	return g.groupDB.CountTotal(ctx, before)
}

func (g *groupDatabase) CreateGroupMember(ctx context.Context, groupMembers []*GroupMemberModel) (err error) {

	//mysql
	if err := g.tx.Transaction(func(tx *gorm.DB) error {
		err := g.groupMember.Create(ctx, groupMembers)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	//redis
	groupMemberIDs := make([]string, len(groupMembers))
	for k, v := range groupMembers {
		groupMemberIDs[k] = strconv.Itoa(int(v.UserID))
	}
	if err := g.groupRDB.AddMemberIDs(ctx, strconv.Itoa(int(groupMembers[0].GroupID)), groupMemberIDs); err != nil {
		return err
	}

	return nil
}

func (g *groupDatabase) DeleteGroupMember(ctx context.Context, groupID uint64, userIDs []uint64) error {
	//mysql
	if err := g.tx.Transaction(func(tx *gorm.DB) error {
		err := g.groupMember.Delete(ctx, groupID, userIDs)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	//redis
	membersID := make([]string, len(userIDs))
	for k, v := range userIDs {
		membersID[k] = strconv.Itoa(int(v))
	}
	g.groupRDB.DelMemberIDs(ctx, strconv.Itoa(int(groupID)), membersID)
	return nil
}

func (g *groupDatabase) TakeGroupMember(ctx context.Context, groupID uint64, userID uint64) (groupMember *GroupMemberModel, err error) {
	return g.groupMember.Take(ctx, groupID, userID)
}
func (g *groupDatabase) FindGroupMemberIDs(ctx context.Context, groupID uint64) (groupMemberIDs []uint64, err error) {
	//redis
	groupMemberIDs, err = g.groupRDB.GetMemberIDs(ctx, strconv.Itoa(int(groupID)))
	if err != nil {
		return nil, err
	}
	if len(groupMemberIDs) == 0 {
		//mysql
		groupMemberIDs, err = g.groupMember.FindGroupMemberIDs(ctx, groupID)
		if err != nil {
			return nil, err
		}
	}
	return groupMemberIDs, nil
}
func (g *groupDatabase) FindUserJoinedGroupID(ctx context.Context, userID uint64) (groupIDs []uint64, err error) {
	return g.groupMember.FindUserJoinedGroupID(ctx, userID)
}

func (g *groupDatabase) TakeGroupMemberNum(ctx context.Context, groupID uint64) (count int64, err error) {
	return g.groupMember.TakeGroupMemberNum(ctx, groupID)
}
func (g *groupDatabase) JoinGroup(ctx context.Context, groupID, inviterID uint64, users []*UserModel) error {
	//user to member
	groupMembers := make([]*GroupMemberModel, len(users))
	for k, user := range users {
		groupMembers[k] = &GroupMemberModel{GroupID: groupID, UserID: user.UserID, NickName: user.NickName, FaceURL: user.FaceURL, JoinTime: time.Now(), InviterUserID: inviterID}
	}
	return g.groupMember.Create(ctx, groupMembers)
}
