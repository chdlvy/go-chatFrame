package orm

import (
	"gorm.io/gorm"
	u "server/utils"
	"strconv"
	"time"
)

type Comment struct {
	CommentId       uint64 `gorm:"comment_id;primaryKey"json:"commentId"`
	UserId          uint64 `gorm:"user_id"json:"userId"`
	VideoId         uint64 `gorm:"video_id"json:"videoId"`
	Content         string `gorm:"content"json:"content"`
	LikeCount       uint64 `gorm:"like_count"json:"likeCount"`
	CreateAt        string `gorm:"create_at"json:"createAt"`
	ParentCommentId uint64 `gorm:"parent_comment_id"json:"parentCommentId"`
}

func AddCommentToSql(userId uint64, videoId uint64, content string, parentCommentId uint64) (commentId *Comment, err error) {
	db := GetSqlConn()
	createAt := strconv.Itoa(int(time.Now().Unix()))
	comment := Comment{UserId: userId, VideoId: videoId, Content: content, CreateAt: createAt, ParentCommentId: parentCommentId}
	err = db.Create(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
func GetChildComment(videoId uint64, parentCommentId uint64, row uint64) ([]Comment, error) {
	db := GetSqlConn()
	var childComment []Comment
	err := db.Limit(int(row)).Where("video_id = ? and parent_comment_id = ?", videoId, parentCommentId).Order("like_count desc").Find(&childComment).Error
	if err != nil {
		return nil, err
	}
	return childComment, nil
}

func GetCommentSql(videoId uint64, page uint64) ([]map[string]interface{}, error) {
	db := GetSqlConn()
	c := []Comment{}
	err := db.Offset(int(page)*20).Limit(20).Where("video_id = ?", videoId).Order("like_count desc").Find(&c).Error
	if err != nil {
		return nil, err
	}
	var arr []map[string]interface{}
	for _, v := range c {
		arr = append(arr, u.StructToMap(v))
	}
	return arr, nil

}

func LikeCommentSql(commentId uint64) error {
	db := GetSqlConn()
	err := db.Model(&Comment{}).Where("comment_id = ?", commentId).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func CancelLikeCommentSql(commentId uint64) error {
	db := GetSqlConn()
	err := db.Model(&Comment{}).Where("comment_id = ?", commentId).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}
