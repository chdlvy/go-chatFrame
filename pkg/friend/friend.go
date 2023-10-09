package friend

import (
	"context"
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/chdlvy/go-chatFrame/pkg/common/db/cache"
	"log"
)

type FriendServer struct {
	friendDatabase db.FriendDatabase
}

func NewFriendServer() *FriendServer {
	dbconn, err := db.NewGormDB()
	if err != nil {
		log.Println(err)
	}
	friendDB := db.NewFriendGorm(dbconn)
	friendRequestDB := db.NewFriendRequestGorm(dbconn)
	rdb, err := cache.NewRedis()
	if err != nil {
		log.Fatal(err)
	}
	friendRdb := cache.NewFriendCache(rdb)
	friendServer := &FriendServer{
		friendDatabase: db.NewFriendDatabase(friendDB, friendRequestDB, friendRdb),
	}
	return friendServer
}

func StartFriendServer() error {

	//mysql
	gdb, err := db.NewGormDB()
	if err != nil {
		return err
	}
	if err = gdb.AutoMigrate(&db.FriendModel{}, &db.FriendRequestModel{}); err != nil {
		return err
	}
	return nil
}
func (f *FriendServer) ApplyAddFriend(fromUserID, toUserID uint64, reqMsg string) error {
	return f.friendDatabase.AddFriendRequest(context.Background(), fromUserID, toUserID, reqMsg)
}
func (f *FriendServer) AgreeFriendReq(fromUserID, toUserID uint64) error {
	return f.friendDatabase.AgreeFriendRequest(context.Background(), fromUserID, toUserID)
}

func (f *FriendServer) RefuseFriendReq(fromUserID, toUserID uint64) error {
	return f.friendDatabase.RefuseFriendRequest(context.Background(), fromUserID, toUserID)
}
func (f *FriendServer) DeleteFriend(ownerUserID, FriendID uint64) error {
	return f.friendDatabase.Delete(context.Background(), ownerUserID, FriendID)
}
func (f *FriendServer) SetFriendRemark(ownerUserID, FriendID uint64, remark string) error {
	return f.friendDatabase.UpdateRemark(context.Background(), ownerUserID, FriendID, remark)
}
func (f *FriendServer) IsFriend(userID1, userID2 uint64) (bool, error) {
	return f.friendDatabase.CheckIn(context.Background(), userID1, userID2)
}

// 获取收到的好友请求
func (f *FriendServer) GetFriendReqToMe(userID uint64, pageNumber, showNumber int32) (friendRequests []*db.FriendRequestModel, err error) {
	return f.friendDatabase.GetFriendReqToMe(context.Background(), userID, pageNumber, showNumber)
}

// 获取发出的好友请求
func (f *FriendServer) GetFriendReqFromMe(userID uint64, pageNumber, showNumber int32) (friendRequests []*db.FriendRequestModel, err error) {
	return f.friendDatabase.GetFriendReqFromMe(context.Background(), userID, pageNumber, showNumber)
}

// 获取所有好友
func (f *FriendServer) GetFriendList(ownerUserID uint64) (friends []*db.FriendModel, err error) {
	return f.friendDatabase.FindFriends(context.Background(), ownerUserID)
}
