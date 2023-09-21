package db

import (
	"context"
	"gorm.io/gorm"
)

type FriendRequestGorm struct {
	*MetaDB
}

func NewFriendRequestGorm(db *gorm.DB) FriendRequestModelInterface {
	return &FriendRequestGorm{NewMetaDB(db, &FriendRequestModel{})}
}

func (f *FriendRequestGorm) Create(ctx context.Context, friendRequests []*FriendRequestModel) (err error) {
	return f.db(ctx).Create(&friendRequests).Error
}
func (f *FriendRequestGorm) Update(ctx context.Context, friendRequest *FriendRequestModel) (err error) {
	return f.db(ctx).Where("from_user_id = ? and to_user_id = ?", friendRequest.FromUserID, friendRequest.ToUserID).Updates(*friendRequest).Error
}

func (f *FriendRequestGorm) Delete(ctx context.Context, fromUserID, toUserID uint64) (err error) {
	return f.db(ctx).Where("from_user_id = ? and to_user_id = ?", fromUserID, toUserID).Delete(FriendRequestModel{}).Error
}
func (f *FriendRequestGorm) Take(ctx context.Context, fromUserID, toUserID uint64) (friendRequest *FriendRequestModel, err error) {
	friendRequest = &FriendRequestModel{}
	return friendRequest, f.db(ctx).Where("from_user_id = ? and to_user_id = ?", fromUserID, toUserID).Take(friendRequest).Error
}

func (f *FriendRequestGorm) FindToUserID(ctx context.Context, toUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error) {
	return friendRequests, f.db(ctx).Where("to_user_id = ?", toUserID).Limit(int(showNumber)).Offset(int(pageNumber-1) * int(showNumber)).Find(&friendRequests).Error
}

func (f *FriendRequestGorm) FindFromUserID(ctx context.Context, fromUserID uint64, pageNumber, showNumber int32) (friendRequests []*FriendRequestModel, err error) {
	return friendRequests, f.db(ctx).Where("to_user_id = ?", fromUserID).Limit(int(showNumber)).Offset(int(pageNumber-1) * int(showNumber)).Find(&friendRequests).Error
}
