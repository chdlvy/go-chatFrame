package group

import (
	"context"
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type GroupServer struct {
	GroupDB db.GroupDatabase
}

func NewGroupServer() *GroupServer {

	return &GroupServer{
		GroupDB: db.InitGroupDatabase(db.DBConn),
	}
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
	g.GroupDB.JoinGroup(context.Background(), groupID, inviterID, []*db.UserModel{user})
	return nil
}

func (g *GroupServer) GenGroupID(groupID *uint64) {
	id := uuid.New()
	gid, _ := strconv.Atoi(id.String())
	*groupID = uint64(gid)
}

func (g *GroupServer) KickMember(groupID, MemberID uint64) error {
	return g.GroupDB.DeleteGroupMember(context.Background(), groupID, []uint64{MemberID})
}
func (g *GroupServer) QuitGroup(ctx context.Context, groupID, userID uint64) error {
	return g.GroupDB.DeleteGroupMember(ctx, groupID, []uint64{userID})
}
