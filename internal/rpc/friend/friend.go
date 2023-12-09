package friend

import (
	"context"
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/chdlvy/go-chatFrame/pkg/common/db/cache"
	friendPb "github.com/chdlvy/protocol/friend"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type friendServer struct {
	friendDatabase db.FriendDatabase
}

func Start() error {
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
	friendServer := &friendServer{
		friendDatabase: db.NewFriendDatabase(friendDB, friendRequestDB, friendRdb),
	}
	//grpcServer := grpc.NewServer()
	//grpcServer.RegisterService()
	return friendServer
}

type FriendRPCServer struct {
	friendPb.UnimplementedFriendServer
}

func (friendRPCServer *FriendRPCServer) ApplyToAddFriend(context.Context, *friendPb.ApplyToAddFriendReq) (*friendPb.ApplyToAddFriendResp, error) {

	return nil, status.Errorf(codes.Unimplemented, "method ApplyToAddFriend not implemented")
}
func (friendRPCServer *FriendRPCServer) GetPaginationFriendsApplyTo(context.Context, *friendPb.GetPaginationFriendsApplyToReq) (*friendPb.GetPaginationFriendsApplyToResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaginationFriendsApplyTo not implemented")
}
func (friendRPCServer *FriendRPCServer) GetPaginationFriendsApplyFrom(context.Context, *friendPb.GetPaginationFriendsApplyFromReq) (*friendPb.GetPaginationFriendsApplyFromResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaginationFriendsApplyFrom not implemented")
}
func (friendRPCServer *FriendRPCServer) AddBlack(context.Context, *friendPb.AddBlackReq) (*friendPb.AddBlackResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBlack not implemented")
}
func (friendRPCServer *FriendRPCServer) RemoveBlack(context.Context, *friendPb.RemoveBlackReq) (*friendPb.RemoveBlackResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveBlack not implemented")
}
func (friendRPCServer *FriendRPCServer) IsFriend(context.Context, *friendPb.IsFriendReq) (*friendPb.IsFriendResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFriend not implemented")
}
func (friendRPCServer *FriendRPCServer) IsBlack(context.Context, *friendPb.IsBlackReq) (*friendPb.IsBlackResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsBlack not implemented")
}
func (friendRPCServer *FriendRPCServer) GetPaginationBlacks(context.Context, *friendPb.GetPaginationBlacksReq) (*friendPb.GetPaginationBlacksResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaginationBlacks not implemented")
}
func (friendRPCServer *FriendRPCServer) DeleteFriend(context.Context, *friendPb.DeleteFriendReq) (*friendPb.DeleteFriendResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFriend not implemented")
}
func (friendRPCServer *FriendRPCServer) RespondFriendApply(context.Context, *friendPb.RespondFriendApplyReq) (*friendPb.RespondFriendApplyResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RespondFriendApply not implemented")
}
func (friendRPCServer *FriendRPCServer) SetFriendRemark(context.Context, *friendPb.SetFriendRemarkReq) (*friendPb.SetFriendRemarkResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetFriendRemark not implemented")
}
func (friendRPCServer *FriendRPCServer) GetPaginationFriends(context.Context, *friendPb.GetPaginationFriendsReq) (*friendPb.GetPaginationFriendsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaginationFriends not implemented")
}
func (friendRPCServer *FriendRPCServer) GetFriendIDs(context.Context, *friendPb.GetFriendIDsReq) (*friendPb.GetFriendIDsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriendIDs not implemented")
}
