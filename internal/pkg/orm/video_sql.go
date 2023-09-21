package orm

import "time"

type Videos struct {
	Id           uint64    `gorm:"id;primaryKey"` //视频id
	WriterId     uint64    `gorm:"writerId"`      //作者id
	VideoSrc     string    `gorm:"videoSrc"`      //视频链接
	Title        string    `gorm:"title" `        //标题
	Description  string    `gorm:"description" `  //描述
	Cover        string    `gorm:"cover" `        //封面
	Duration     time.Time `gorm:"duration" `     //时长
	PlayCount    uint64    `gorm:"playCount" `    //播放量
	LikeCount    uint64    `gorm:"likeCount" `    //点赞数
	CollectCount uint64    `gorm:"collectCount" ` //收藏数
	ShareCount   uint64    `gorm:"shareCount" `   //分享数
	CommentCount uint64    `gorm:"commentCount" ` //评论数
	IsLike       uint8     `gorm:"isLike"`        //检查用户是否点赞过
}

// 收藏表
type Collect struct {
	Uid     uint64 `gorm:"uid"`
	VideoId uint64 `gorm:"videoId"`
}

func AddVideoToSql(v map[string]interface{}) error {
	db := GetSqlConn()
	video := Videos{
		WriterId:    v["writerId"].(uint64),
		VideoSrc:    v["videoSrc"].(string),
		Title:       v["title"].(string),
		Description: v["description"].(string),
		Cover:       v["cover"].(string),
		Duration:    v["duration"].(time.Time),
	}
	err := db.Create(&video).Error
	if err != nil {
		return err
	}
	return nil
}
func DelVideo(videoId uint64) error {
	if err := db.Delete(&Videos{}, videoId).Error; err != nil {
		return err
	}
	return nil
}

// 收藏视频
func CollectVideo(uid uint64, videoId uint64) error {
	collect := &Collect{
		Uid:     uid,
		VideoId: videoId,
	}
	err := db.Create(&collect).Error
	if err != nil {
		return err
	}
	return nil
}

func UnCollectVideo(uid uint64, videoId uint64) error {
	result := db.Where("uid = ? and videoId = ?", uid, videoId).Delete(&Collect{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetVideos(page uint64, rowsNum uint64) []Videos {
	var video []Videos
	//推荐算法
	db.Find(&video)

	return video
}
