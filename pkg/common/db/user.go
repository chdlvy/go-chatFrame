package db

import (
	"context"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type UserModel struct {
	UserID     uint64    `gorm:"column:user_id;primary_key;size:64"`
	NickName   string    `gorm:"column:name;size:255"`
	FaceURL    string    `gorm:"column:face_url;size:255"`
	Token      string    `gorm:"colum:token"`
	CreateTime time.Time `gorm:"column:create_time;index:create_time;autoCreateTime"`
}
type UserModelInterface interface {
	Create(ctx context.Context, users []*UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, user *UserModel) (err error)
	// 获取指定用户信息  不存在，也不返回错误
	Find(ctx context.Context, userIDs []uint64) (users []*UserModel, err error)
	// 获取某个用户信息  不存在，则返回错误
	Take(ctx context.Context, userID uint64) (user *UserModel, err error)
	GetAllUserID(ctx context.Context, pageNumber, showNumber int32) (userIDs []uint64, err error)
	// 获取用户总数
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
}

func StartUserServer() error {
	db, err := NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		return nil
	}
	return nil
}
func GenUserID() uint64 {
	id := uuid.New()
	uid, _ := strconv.Atoi(id.String())
	return uint64(uid)
}
