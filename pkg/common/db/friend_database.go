package db

import (
	"context"
	"errors"
	"github.com/chdlvy/go-chatFrame/constant"
	"github.com/chdlvy/go-chatFrame/pkg/common/db/cache"
	"gorm.io/gorm"
	"log"
	"time"
)

type FriendDatabase interface {
	//检查friendID是否在好友列表中
	CheckIn(ctx context.Context, ownUserID, friendID uint64) (bool, error)
	// 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID uint64, reqMsg string) (err error)
	// 成为好友，这里不执行设置备注操作
	BecomeFriends(ctx context.Context, ownerUserID uint64, friendUserID uint64, addSource int32) (err error)
	// 拒绝好友申请
	RefuseFriendRequest(ctx context.Context, fromUserID, toUserID uint64) (err error)
	// 同意好友申请
	AgreeFriendRequest(ctx context.Context, fromUserID, toUserID uint64) (err error)
	// 删除好友
	Delete(ctx context.Context, ownerUserID uint64, friendUserID uint64) (err error)
	// 更新好友备注
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID uint64, remark string) (err error)
	//获取收到的好友请求
	GetFriendReqToMe(ctx context.Context, toUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error)
	// 获取发出的好友请求
	GetFriendReqFromMe(ctx context.Context, fromUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error)
	//获取所有好友
	FindFriends(ctx context.Context, ownerUserID uint64) (friends []*FriendModel, err error)
}
type friendDatabase struct {
	friend        FriendModelInterface
	friendRequest FriendRequestModelInterface
	friendRdb     cache.FriendCache
	tx            *gorm.DB
}

func (f *friendDatabase) FindFriends(ctx context.Context, ownerUserID uint64) (friends []*FriendModel, err error) {
	return f.friend.FindFriends(ctx, ownerUserID)
}

func NewFriendDatabase(friend FriendModelInterface, friendRequest FriendRequestModelInterface, rdb cache.FriendCache) FriendDatabase {
	return &friendDatabase{friend: friend, friendRequest: friendRequest, friendRdb: rdb, tx: DBConn}
}
func (f *friendDatabase) CheckIn(ctx context.Context, ownUserID, friendID uint64) (bool, error) {
	friendIDs, err := f.friend.FindFriendIDs(ctx, ownUserID)
	if err != nil {
		return false, err
	}
	for _, v := range friendIDs {
		if v == friendID {
			return true, nil
		}
	}
	return false, nil
}

func (f *friendDatabase) AddFriendRequest(ctx context.Context, fromUserID, toUserID uint64, reqMsg string) (err error) {
	fr, err := f.friendRequest.Take(ctx, fromUserID, toUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == nil && fr.HandleResult == -1 {
		fr.HandleResult = constant.FriendResponseNotHandle
		fr.HandleTime = time.Now()
		fr.ReqMsg = reqMsg
		if err = f.friendRequest.Update(ctx, fr); err != nil {
			return err
		}
		return
	}

	//如果已经存在请求记录，则不添加
	if err == nil {
		return
	}
	friendRequest := &FriendRequestModel{FromUserID: fromUserID, ToUserID: toUserID, ReqMsg: reqMsg, CreateTime: time.Now(), HandleTime: time.Unix(0, 0)}
	if err := f.friendRequest.Create(ctx, []*FriendRequestModel{friendRequest}); err != nil {
		return err
	}
	return nil

}

func (f *friendDatabase) BecomeFriends(ctx context.Context, ownerUserID uint64, friendUserID uint64, addSource int32) (err error) {
	//要添加两条记录，因为互为好友
	// 先判断是否在好友表，如果在则不插入
	if err := f.tx.Transaction(func(tx *gorm.DB) error {
		hasFriend, err := f.CheckIn(ctx, ownerUserID, friendUserID)
		if err != nil {
			log.Println(err)
			return err
		}
		//如果不在好友表，则添加好友
		if !hasFriend {
			friend := &FriendModel{OwnerUserID: ownerUserID, FriendUserID: friendUserID, AddSource: addSource, CreateTime: time.Now()}
			f.friend.Create(ctx, []*FriendModel{friend})
		}

		//对对方的好友操作
		hasFriend2, err := f.CheckIn(ctx, friendUserID, ownerUserID)
		if err != nil {
			log.Println(err)
			return err
		}
		//如果不在好友表，则添加好友
		if !hasFriend2 {
			friend := &FriendModel{OwnerUserID: friendUserID, FriendUserID: ownerUserID, AddSource: addSource, CreateTime: time.Now()}
			f.friend.Create(ctx, []*FriendModel{friend})
		}
		return nil
	}); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (f *friendDatabase) RefuseFriendRequest(ctx context.Context, fromUserID, toUserID uint64) (err error) {
	friendRequest, err := f.friendRequest.Take(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}

	if friendRequest.HandleResult != 0 {
		return errors.New("this friend request has been handle")
	}
	friendRequest.HandleResult = constant.FriendResponseRefuse
	friendRequest.HandleTime = time.Now()
	return f.friendRequest.Update(ctx, friendRequest)
}

func (f *friendDatabase) AgreeFriendRequest(ctx context.Context, fromUserID, toUserID uint64) (err error) {
	friendRequest, err := f.friendRequest.Take(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}

	if friendRequest.HandleResult != 0 {
		return errors.New("this friend request has been handle")
	}
	friendRequest.HandleResult = constant.FriendResponseAgree
	friendRequest.HandleTime = time.Now()
	if err := f.friendRequest.Update(ctx, friendRequest); err != nil {
		return err
	}
	//成为好友
	return f.BecomeFriends(ctx, fromUserID, toUserID, constant.BecomeFriendByApply)
}

func (f *friendDatabase) Delete(ctx context.Context, ownerUserID uint64, friendUserID uint64) (err error) {
	return f.friend.Delete(ctx, ownerUserID, []uint64{friendUserID})
}
func (f *friendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID uint64, remark string) (err error) {
	return f.friend.UpdateRemark(ctx, ownerUserID, friendUserID, remark)
}
func (f *friendDatabase) GetFriendReqToMe(ctx context.Context, toUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error) {
	return f.friendRequest.FindToUserID(ctx, toUserID, pageNumber, showNumber)
}

func (f *friendDatabase) GetFriendReqFromMe(ctx context.Context, fromUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error) {
	return f.friendRequest.FindFromUserID(ctx, fromUserID, pageNumber, showNumber)
}
