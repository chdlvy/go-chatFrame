package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	sql "server/internal/pkg/orm"
	rd "server/internal/pkg/redis"
	t "server/jwt"
	"server/middleware"
	util "server/utils"
	"strconv"
	"time"
)

type Routes struct {
}

func (this *Routes) InitRouter(r *gin.Engine) {
	r.Use(middleware.Cors())

	r.POST("login", this.login)
	r.POST("publish", this.publish)
	r.POST("addComment", this.addComment)
	r.GET("getComment", this.getComment)
	r.GET("getChildComment", this.getChildComment)

	friendRouterGroup := r.Group("/friend")
	{
		friendRouterGroup.POST("add_friend")

	}
}

func (this *Routes) login(c *gin.Context) {
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	err := sql.IsExistUser(phone)
	if err != nil {
		// 注册账号
		if err.Error() == "账号未注册" {
			var newUser = &sql.Users{
				Phone:        phone,
				UserName:     "用户" + strconv.Itoa(int(time.Now().Unix())),
				Password:     password,
				Avatar:       "defaultImg",
				Gender:       2,
				Birthday:     "",
				PublishCount: 0,
				LikedCount:   0,
				CreateTime:   time.Now(),
			}

			//添加user
			userId, err := sql.SignUpUser(newUser)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg": err.Error(),
				})
			} else {
				// 注册完毕直接登录
				m := map[string]interface{}{
					"phone":    phone,
					"password": password,
				}
				token := t.GetToken(m)
				//获取user信息
				userInfo, _ := sql.GetUserFromId(userId)

				c.JSON(http.StatusOK, gin.H{
					"status": 1,
					"msg":    "注册并登录成功",
					"token":  token,
					"user":   userInfo,
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		}
	} else {
		// 登录账号并返回token
		m := map[string]interface{}{
			"phone":    phone,
			"password": password,
		}
		token := t.GetToken(m)
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "登录成功",
			"token":  token,
		})
	}
}

func (this *Routes) publish(c *gin.Context) {
	var video sql.Videos
	c.ShouldBind(&video)
	video_map := util.StructToMap(video)
	err := sql.AddVideoToSql(video_map)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": "发布失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": "发布成功",
	})
}

func (this *Routes) delVideo(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.PostForm("id"))
	err := sql.DelVideo(uint64(videoId))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": "删除失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": "删除成功",
		})
	}
}
func (this *Routes) addComment(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.PostForm("videoId"))
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	parentCommentId, _ := strconv.Atoi(c.PostForm("commentId"))
	content := c.PostForm("content")
	var err error
	//检查videoId和commentId是否存在=============(need to do)

	if parentCommentId != 0 {
		fmt.Println("setChildComment")
		err = rd.SetChildComment(uint64(parentCommentId), uint64(userId), uint64(videoId), content)
	} else {
		fmt.Println("setComment")
		err = rd.SetComment(uint64(userId), uint64(videoId), content)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"data":  "发布评论失败",
			"error": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": "发布评论成功",
	})
}
func (this *Routes) getComment(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.Query("videoId"))
	list, err := rd.GetCommentFromVideoId(uint64(videoId), 10)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
	}
	if len(list) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": "暂时没有人评论",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"list": list,
	})
}

func (this *Routes) getChildComment(c *gin.Context) {
	videoId, _ := strconv.Atoi(c.Query("videoId"))
	commentId, _ := strconv.Atoi(c.Query("commentId"))
	rd.GetChildComment(uint64(commentId), uint64(videoId))
}
