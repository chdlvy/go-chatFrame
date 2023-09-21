package db

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// 操作数据库来执行group的操作
type GroupGorm struct {
	*MetaDB
}

var _ GroupModelInterface = (*GroupGorm)(nil)

func NewGroupDB(db *gorm.DB) GroupModelInterface {
	return &GroupGorm{NewMetaDB(db, &GroupModel{})}
}

func (g *GroupGorm) Create(ctx context.Context, group *GroupModel) (err error) {
	return g.DB.Create(&group).Error
}

func (g *GroupGorm) Take(ctx context.Context, groupID uint64) (group *GroupModel, err error) {
	group = &GroupModel{}
	return group, g.DB.Where("group_id", groupID).Take(group).Error
}
func (g *GroupGorm) Find(ctx context.Context, groupIDs []uint64) (groups []*GroupModel, err error) {
	return groups, g.DB.Where("group_id in ?", groupIDs).Find(&groups).Error
}

func (g *GroupGorm) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	db := g.db(ctx)
	if before != nil {
		db = db.Where("create_time < ?", before)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
