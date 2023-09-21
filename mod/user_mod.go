package mod

import "time"

type UserMod struct {
	Id           uint64    `redis:"id"`
	UserName     string    `redis:"username" form:"userName"`
	Password     string    `redis:"password" form:"password"`
	Phone        string    `redis:"phone" form:"phone"`
	Avatar       string    `redis:"avatar" form:"avatar"`
	Gender       uint8     `redis:"gender" form:"gender"`
	Birthday     string    `redis:"birthday" form:"birthday"`
	PublishCount uint64    `redis:"publish_count" form:"publishCount"`
	LikedCount   uint64    `redis:"liked_count" form:"likedCount"`
	CreateTime   time.Time `redis:"create_time" form:"createTime"`
}
