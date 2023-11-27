package msggateway

import (
	"context"
	"github.com/chdlvy/go-chatFrame/pkg/chatlog"
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/chdlvy/go-chatFrame/pkg/friend"
	"github.com/chdlvy/go-chatFrame/pkg/group"
	"log"
)

type MidServer struct {
	GroupServer   *group.GroupServer
	UserServer    db.UserModelInterface
	FriendServer  *friend.FriendServer
	ChatLogServer *chatlog.ChatLogServer
}

func NewMidServer() *MidServer {
	gdb, err := db.NewGormDB()
	if err != nil {
		log.Fatal(err)
	}
	Userdb := db.NewUserGorm(gdb)
	Gs := group.NewGroupServer()
	Fs := friend.NewFriendServer()
	Clog := chatlog.NewChatLogServer()
	return &MidServer{GroupServer: Gs, UserServer: Userdb, FriendServer: Fs, ChatLogServer: Clog}
}

func StartMidServer() error {
	if err := group.StartGroupServer(); err != nil {
		return err
	}
	if err := db.StartUserServer(); err != nil {
		return err
	}
	if err := friend.StartFriendServer(); err != nil {
		return err
	}
	if err := chatlog.StartChatLogServer(); err != nil {
		return err
	}
	return nil
}
func (m *MidServer) QuitGroup(ctx context.Context, userID, groupID uint64) error {
	//退群
	if err := m.GroupServer.QuitGroup(ctx, groupID, userID); err != nil {
		return err
	}
	//删除消息记录
	if err := m.ChatLogServer.DelGroupAllChatLogWithUser(userID, groupID); err != nil {
		return err
	}
	return nil
}
func (m *MidServer) KickMember(ctx context.Context, userID, groupID uint64) error {
	//退群
	if err := m.GroupServer.QuitGroup(ctx, groupID, userID); err != nil {
		return err
	}
	//删除消息记录
	if err := m.ChatLogServer.DelGroupAllChatLogWithUser(userID, groupID); err != nil {
		return err
	}
	return nil
}
