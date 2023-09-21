package db

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserGorm struct {
	*MetaDB
}

func NewUserGorm(db *gorm.DB) UserModelInterface {
	return &UserGorm{NewMetaDB(db, &UserModel{})}
}

// 插入多条记录
func (u *UserGorm) Create(ctx context.Context, users []*UserModel) (err error) {
	return u.db(ctx).Create(&users).Error
}

// 更新用户信息,零值
func (u *UserGorm) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return u.db(ctx).Where("user_id = ?", userID).Updates(args).Error
}

// 更新多个用户信息 非零值.
func (u *UserGorm) Update(ctx context.Context, user *UserModel) (err error) {
	return u.db(ctx).Updates(user).Error
}

// 获取指定用户信息  不存在，也不返回错误.
func (u *UserGorm) Find(ctx context.Context, userIDs []uint64) (users []*UserModel, err error) {
	err = u.db(ctx).Where("user_id in (?)", userIDs).Find(&users).Error
	return users, err
}

// 获取某个用户信息  不存在，则返回错误.
func (u *UserGorm) Take(ctx context.Context, userID uint64) (user *UserModel, err error) {
	user = &UserModel{}
	err = u.db(ctx).Where("user_id = ?", userID).Take(&user).Error
	return user, err
}

// 获取所有用户ID.
func (u *UserGorm) GetAllUserID(ctx context.Context, pageNumber, showNumber int32) (userIDs []uint64, err error) {
	//查询全部
	if pageNumber == 0 || showNumber == 0 {
		return userIDs, u.db(ctx).Pluck("user_id", &userIDs).Error
	} else {
		//按照范围查询
		return userIDs, u.db(ctx).Limit(int(showNumber)).Offset(int((pageNumber-1)*showNumber)).Pluck("user_id", &userIDs).Error
	}
}

// 查询某个时间点前的用户总数
func (u *UserGorm) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	db := u.db(ctx).Model(&UserModel{})
	if before != nil {
		db = db.Where("create_time < ?", before)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
