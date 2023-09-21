package orm

import (
	"errors"
	"time"
)

type Users struct {
	Id           uint64    `gorm:"id,primaryKey"`
	Password     string    `gorm:"password"`
	UserName     string    `gorm:"userName"`
	Phone        string    `gorm:"phone"`
	Avatar       string    `gorm:"avatar"`
	Gender       uint8     `gorm:"gender"`
	Birthday     string    `gorm:"birthday"`
	PublishCount uint64    `gorm:"publishCount"`
	LikedCount   uint64    `gorm:"likeCount"`
	CreateTime   time.Time `gorm:"createTime"`
}

// 检查用户是否存在
func IsExistUser(phone string) error {
	var user Users
	if err := db.Where("phone = ?", phone).Find(&user).Error; err != nil {
		return err
	}
	if user.Id == 0 {
		return errors.New("账号未注册")
	}
	return nil
}

// 添加user
func SignUpUser(user *Users) (uint64, error) {
	if err := db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.Id, nil
}

func GetUserFromId(id uint64) (Users, error) {
	var user Users
	if err := db.Where("id = ?", id).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
func DestoryUser(id uint64) error {
	if err := db.Delete(&Users{}, id).Error; err != nil {
		return err
	}
	return nil
}
