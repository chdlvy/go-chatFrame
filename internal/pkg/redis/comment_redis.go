package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	sql "server/internal/pkg/orm"
	"strconv"
)

type Comment struct {
}

func SetComment(userId uint64, videoId uint64, content string) error {
	rdb := GetredisConn()
	defer rdb.Close()
	//先插入评论数据到数据库
	commentData, err := sql.AddCommentToSql(userId, videoId, content, 0)
	if err != nil {
		log.Fatal("add comment to mysql err：", err)
		return err
	}
	//同步该评论到redis中
	list := "videoCommentList:" + strconv.Itoa(int(videoId))
	rdb.Do("lpush", list, commentData.CommentId)
	commentStr, _ := json.Marshal(&commentData)
	rdb.Do("hset", "videoId:"+strconv.Itoa(int(videoId))+":comments", commentData.CommentId, commentStr)
	return nil
}

func GetCommentFromVideoId(videoId uint64, row uint64) ([]sql.Comment, error) {
	rdb := GetredisConn()
	defer rdb.Close()
	//rdb.Send("select", "2")
	//先拿到视频对应的评论id
	//videoCommentList:videoid保存视频的所有评论id
	list := "videoCommentList:" + strconv.Itoa(int(videoId))
	comments, err := redis.Strings(rdb.Do("lrange", list, 0, row))
	if err != nil {
		return nil, err
	}
	//不存在该videoId
	if len(comments) == 0 {
		return nil, errors.New("not exist videoId")
	}
	//fmt.Println(reflect.TypeOf(comments[0]))
	//根据评论id拿到具体的评论数据
	commentArr := make([]sql.Comment, len(comments))
	commentMap := "videoId:" + strconv.Itoa(int(videoId)) + ":comments"
	for k, v := range comments {
		commentStr, _ := redis.String(rdb.Do("hget", commentMap, v))

		var commentData sql.Comment
		json.Unmarshal([]byte(commentStr), &commentData)
		//构造[]comment
		commentArr[k] = commentData
	}

	return commentArr, nil
}

func SetChildComment(commentId uint64, userId uint64, videoId uint64, content string) error {
	rdb := GetredisConn()
	defer rdb.Close()
	//先对set到数据库，然后再同步到redis
	commentData, err := sql.AddCommentToSql(userId, videoId, content, commentId)
	if err != nil {
		log.Fatal("add comment to mysql err：", err)
		return err
	}
	//同步该评论到redis中
	list := "parent:" + strconv.Itoa(int(commentId)) + ":child:" + strconv.Itoa(int(commentData.CommentId))
	rdb.Do("lpush", list, commentData.CommentId)
	commentStr, _ := json.Marshal(&commentData)
	rdb.Do("hset", "videoId:"+strconv.Itoa(int(videoId))+":comments", commentData.CommentId, commentStr)
	return nil
}

// 先从redis中拿，没有再从sql中拿
func GetChildComment(parentCommentId uint64, videoId uint64) {
	childComments, _ := sql.GetChildComment(videoId, parentCommentId, 10)
	fmt.Println(childComments)
	//return
}

func ReplyComment(parentCommentId uint64) {

}
