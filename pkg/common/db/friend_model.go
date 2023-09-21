package db

import (
	"context"
	"gorm.io/gorm"
)

type FriendGorm struct {
	*MetaDB
}

func NewFriendGorm(db *gorm.DB) FriendModelInterface {
	return &FriendGorm{NewMetaDB(db, &FriendModel{})}
}
func (f *FriendGorm) Create(ctx context.Context, friends []*FriendModel) (err error) {
	return f.db(ctx).Create(&friends).Error
}

func (f *FriendGorm) Delete(ctx context.Context, ownerUserID uint64, friendUserIDs []uint64) (err error) {
	return f.db(ctx).Where("owner_user_id = ? and friend_user_id in (?)", ownerUserID, friendUserIDs).Delete(&FriendModel{}).Error
}

func (f *FriendGorm) Update(ctx context.Context, ownerUserID, friendUserID uint64, info map[string]interface{}) (err error) {
	return f.db(ctx).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Updates(info).Error
}

func (f *FriendGorm) UpdateRemark(ctx context.Context, ownerUserID, friendUserID uint64, remark string) (err error) {
	return f.db(ctx).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Update("remark", remark).Error
}

func (f *FriendGorm) Take(ctx context.Context, ownerUserID, friendUserID uint64) (friend *FriendModel, err error) {
	friend = &FriendModel{}
	return friend, f.db(ctx).Where("owner_user_id = ? and friend_user_id = ?", ownerUserID, friendUserID).Take(&friend).Error
}

func (f *FriendGorm) FindFriends(ctx context.Context, ownerUserID uint64) (friends []*FriendModel, err error) {
	friendIDs, err := f.FindFriendIDs(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	return friends, f.db(ctx).Where("owner_user_id = ? and friend_user_id in (?)", ownerUserID, friendIDs).Find(&friends).Error
}

func (f *FriendGorm) FindFriendIDs(ctx context.Context, ownerUserID uint64) (friendIDs []uint64, err error) {

	return friendIDs, f.db(ctx).Where("owner_user_id = ?", ownerUserID).Pluck("friend_user_id", &friendIDs).Error
}
