package group

import (
	"context"
	"errors"
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"math/rand"
	"time"
)

type GroupServer struct {
	GroupIDs map[int]bool
	Groups   map[uint64]*Group
	GroupDB  db.GroupDatabase
}

func NewGroupServer() *GroupServer {

	return &GroupServer{
		Groups:   make(map[uint64]*Group),
		GroupIDs: make(map[int]bool),
		GroupDB:  db.InitGroupDatabase(db.DBConn),
	}
}
func (g *GroupServer) GetGroup(groupID uint64) (*Group, error) {
	v, ok := g.Groups[groupID]
	if !ok {
		return nil, errors.New("group is not exist")
	}
	return v, nil
}

func (g *GroupServer) CreateGroup(groupInfo *GroupInfo, userIDs []uint64) (*Group, error) {
	userDB := db.NewUserGorm(db.DBConn)
	creator, err := userDB.Take(context.Background(), userIDs[0])
	if err != nil {
		return nil, err
	}
	groupID := groupInfo.GroupID
	if groupID == 0 {
		g.GenGroupID(&groupID)
	}
	group := &Group{
		GroupID:       groupID,
		GroupName:     groupInfo.GroupName,
		FaceURL:       groupInfo.FaceURL,
		CreatorUserID: creator.UserID,
		GroupType:     groupInfo.GroupType,
		GroupMember:   make([]*GroupMember, 1),
	}
	dbgroup := &db.GroupModel{
		GroupName:     groupInfo.GroupName,
		FaceURL:       groupInfo.FaceURL,
		CreatorUserID: creator.UserID,
		GroupType:     groupInfo.GroupType,
		CreateTime:    time.Now(),
		GroupID:       groupID,
	}
	groupMember := User2GroupMember(creator)
	groupMember.GroupID = groupID
	groupMember.JoinTime = time.Now()
	groupMember.InviterUserID = 0
	//mysql,container create group and take member to group
	g.GroupDB.CreateGroup(context.Background(), dbgroup, []*db.GroupMemberModel{groupMember})

	g.Groups[groupID] = group
	//creator join group and the other member join group
	err = g.JoinGroup(groupID, creator.UserID, creator.UserID)
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(userIDs); i++ {
		g.JoinGroup(groupID, creator.UserID, userIDs[i])
	}
	return group, nil
}

func (g *GroupServer) JoinGroup(groupID, inviterID, userID uint64) error {
	if _, err := g.GroupDB.TakeGroup(context.Background(), groupID); err != nil {
		return err
	}

	userDB := db.NewUserGorm(db.DBConn)
	user, err := userDB.Take(context.Background(), userID)
	if err != nil {
		return err
	}
	//gm := &GroupMember{
	//	UserID:        user.UserID,
	//	NickName:      user.NickName,
	//	FaceURL:       user.FaceURL,
	//	JoinTime:      time.Now(),
	//	InviterUserID: inviterID,
	//}
	g.GroupDB.JoinGroup(context.Background(), groupID, inviterID, []*db.UserModel{user})
	//g.Groups[groupID].GroupMember = append(g.Groups[groupID].GroupMember, gm)
	return nil
}

func (g *GroupServer) GenGroupID(groupID *uint64) {
	rand.Seed(time.Now().UnixNano())
	for {
		// 生成随机数作为群号，范围在 10000 到 100000 之间
		gID := rand.Intn(90001) + 10000

		// 检查群号是否已被使用
		if !g.GroupIDs[gID] {
			// 群号未被使用，分配该群号
			g.GroupIDs[gID] = true
			*groupID = uint64(gID)
			break
		}
	}
}

func (g *GroupServer) KickMember(groupID, MemberID uint64) error {
	return g.GroupDB.DeleteGroupMember(context.Background(), groupID, []uint64{MemberID})
}
func (g *GroupServer) QuitGroup(ctx context.Context, groupID, userID uint64) error {
	return g.GroupDB.DeleteGroupMember(ctx, groupID, []uint64{userID})
}
