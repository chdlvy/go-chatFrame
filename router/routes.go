package router

import (
	"context"
	token "github.com/chdlvy/go-chatFrame/jwt"
	"github.com/chdlvy/go-chatFrame/middleware"
	"github.com/chdlvy/go-chatFrame/msggateway"
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/chdlvy/go-chatFrame/pkg/group"
	"github.com/chdlvy/go-chatFrame/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Route struct {
	midServer *msggateway.MidServer
}

func (r *Route) InitRouter(e *gin.Engine) {
	userGroup := e.Group("/user")
	userGroup.POST("/login")
	userGroup.POST("/register")
	e.Use(middleware.CheckAuth())
	userGroup.POST("/logout")
	e.Use(middleware.Cors())
	r.midServer = msggateway.NewMidServer()
	e.POST("createGroup", r.createGroup)
	e.POST("createUser", r.createUser)
	e.POST("joinGroup", r.joinGroup)
	e.POST("kickMember", r.kickMember)
	e.POST("quitGroup", r.quitGroup)
	e.POST("applyFriendReq", r.applyFriendReq)
	e.POST("agreeFriendReq", r.agreeFriendReq)
	e.POST("refuseFriendReq", r.refuseFriendReq)
	e.GET("getRecvFriendReq", r.getRecvFriendReq)
	e.GET("getSendFriendReq", r.getSendFriendReq)
	e.POST("deleteFriend", r.deleteFriend)
	e.POST("setFriendRemark", r.setFriendRemark)
	e.GET("getFriendList", r.getFriendList)
}
func (r *Route) createUser(c *gin.Context) {
	user := &db.UserModel{}
	c.ShouldBind(user)
	user.UserID = db.GenUserID()
	user.CreateTime = time.Now()
	user.Token = token.GetToken(utils.StructToMap(user))
	if err := r.midServer.UserServer.Create(context.Background(), []*db.UserModel{user}); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"status": "create failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "create success",
		"user":   user,
	})
}
func (r *Route) createGroup(c *gin.Context) {
	ginfo := &group.GroupInfo{}
	c.ShouldBind(ginfo)
	ginfo.CreateTime = time.Now().Unix()
	ginfo.MemberCount = 1
	if _, err := r.midServer.GroupServer.CreateGroup(ginfo, []uint64{ginfo.CreatorUserID}); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"status": "create group failed",
			"error":  err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "create success",
	})
}
func (r *Route) joinGroup(c *gin.Context) {
	type join struct {
		GroupID   uint64 `json:"groupID,omitempty"`
		InviterID uint64 `json:"inviterID,omitempty"`
		UserID    uint64 `json:"userID,omitempty"`
	}
	joinInfo := &join{}
	c.ShouldBind(joinInfo)
	if err := r.midServer.GroupServer.JoinGroup(joinInfo.GroupID, joinInfo.InviterID, joinInfo.UserID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "join group failed",
			"error":  err,
		})
	}
}
func (r *Route) kickMember(c *gin.Context) {
	type kick struct {
		UserID  uint64 `json:"userID,omitempty"`
		GroupID uint64 `json:"groupID,omitempty"`
	}
	kickInfo := &kick{}
	c.ShouldBind(kickInfo)
	if err := r.midServer.GroupServer.KickMember(kickInfo.GroupID, kickInfo.UserID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "kick user failed",
			"error":  err,
		})
	}
}
func (r *Route) quitGroup(c *gin.Context) {
	type quit struct {
		UserID  uint64 `json:"userID,omitempty"`
		GroupID uint64 `json:"groupID,omitempty"`
	}
	quitInfo := &quit{}
	c.ShouldBind(quitInfo)
	if err := r.midServer.QuitGroup(context.Background(), quitInfo.UserID, quitInfo.GroupID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "quit group failed",
			"error":  err,
		})
	}
}
func (r *Route) applyFriendReq(c *gin.Context) {
	type apply struct {
		FromUserID uint64 `json:"fromUserID,omitempty"`
		ToUserID   uint64 `json:"toUserID,omitempty"`
		ReqMsg     string `json:"reqMsg,omitempty"`
	}
	applyInfo := &apply{}
	c.ShouldBind(applyInfo)
	if err := r.midServer.FriendServer.ApplyAddFriend(applyInfo.FromUserID, applyInfo.ToUserID, applyInfo.ReqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "applyFriend failed",
			"error":  err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
func (r *Route) agreeFriendReq(c *gin.Context) {
	type agree struct {
		FromUserID uint64 `json:"fromUserID,omitempty"`
		ToUserID   uint64 `json:"toUserID,omitempty"`
	}
	agreeInfo := &agree{}
	c.ShouldBind(agreeInfo)
	if err := r.midServer.FriendServer.AgreeFriendReq(agreeInfo.FromUserID, agreeInfo.ToUserID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "agreeFriend failed",
			"error":  err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
func (r *Route) refuseFriendReq(c *gin.Context) {
	type refuse struct {
		FromUserID uint64 `json:"fromUserID,omitempty"`
		ToUserID   uint64 `json:"toUserID,omitempty"`
	}
	refuseInfo := &refuse{}
	c.ShouldBind(refuseInfo)
	if err := r.midServer.FriendServer.RefuseFriendReq(refuseInfo.FromUserID, refuseInfo.ToUserID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "refuseFriend failed",
			"error":  err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}

func (r *Route) getRecvFriendReq(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("userID"))
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))
	showNumber, _ := strconv.Atoi(c.Query("showNumber"))
	err, friendRequests := r.midServer.FriendServer.GetFriendReqToMe(uint64(userID), int32(pageNumber), int32(showNumber))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   friendRequests,
	})

}

func (r *Route) getSendFriendReq(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("userID"))
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))
	showNumber, _ := strconv.Atoi(c.Query("showNumber"))
	err, friendRequests := r.midServer.FriendServer.GetFriendReqFromMe(uint64(userID), int32(pageNumber), int32(showNumber))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   friendRequests,
	})

}
func (r *Route) deleteFriend(c *gin.Context) {
	type req struct {
		UserID   uint64
		FriendID uint64
	}
	reqInfo := &req{}
	c.ShouldBind(reqInfo)
	if err := r.midServer.FriendServer.DeleteFriend(reqInfo.UserID, reqInfo.FriendID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
func (r *Route) getFriendList(c *gin.Context) {
	ownerUserID, _ := strconv.Atoi(c.Query("userID"))
	friends, err := r.midServer.FriendServer.GetFriendList(uint64(ownerUserID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   friends,
	})
}
func (r *Route) setFriendRemark(c *gin.Context) {
	type req struct {
		UserID   uint64 `json:"userID,omitempty"`
		FriendID uint64 `json:"friendID,omitempty"`
		remark   string `json:"remark,omitempty"`
	}
	reqInfo := &req{}
	c.ShouldBind(reqInfo)
	if err := r.midServer.FriendServer.SetFriendRemark(reqInfo.UserID, reqInfo.FriendID, reqInfo.remark); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
