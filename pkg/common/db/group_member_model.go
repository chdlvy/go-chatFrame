package db

import (
	"context"
	"gorm.io/gorm"
)

type GroupMemberGorm struct {
	*MetaDB
}

var _ GroupMemberModelInterface = (*GroupMemberGorm)(nil)

func NewGroupMemberDB(db *gorm.DB) GroupMemberModelInterface {
	return &GroupMemberGorm{NewMetaDB(db, &GroupMemberModel{})}
}

func (g *GroupMemberGorm) Create(ctx context.Context, groupMembers []*GroupMemberModel) (err error) {
	return g.db(ctx).Create(&groupMembers).Error
}
func (g *GroupMemberGorm) Delete(ctx context.Context, groupID uint64, userIDs []uint64) (err error) {
	return g.db(ctx).Where("group_id = ? and user_id in (?)", groupID, userIDs).Delete(&GroupMemberModel{}).Error
}

func (g *GroupMemberGorm) Take(ctx context.Context, groupID uint64, userID uint64) (groupMember *GroupMemberModel, err error) {
	groupMember = &GroupMemberModel{}
	return groupMember, g.db(ctx).Where("group_id = ? and user_id = ?", groupID, userID).Take(groupMember).Error
}
func (g *GroupMemberGorm) FindGroupMemberIDs(ctx context.Context, groupID uint64) (groupMemberIDs []uint64, err error) {
	return groupMemberIDs, g.db(ctx).Where("group_id = ?", groupID).Pluck("user_id", &groupMemberIDs).Error
}
func (g *GroupMemberGorm) FindUserJoinedGroupID(ctx context.Context, userID uint64) (groupIDs []uint64, err error) {

	return groupIDs, g.db(ctx).Where("user_id = ?", userID).Pluck("group_id", &groupIDs).Error
}

func (g *GroupMemberGorm) TakeGroupMemberNum(ctx context.Context, groupID uint64) (count int64, err error) {
	return count, g.db(ctx).Where("group_id = ?", groupID).Count(&count).Error
}
