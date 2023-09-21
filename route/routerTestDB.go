package route

import (
	token "chatFrame/jwt"
	"chatFrame/middleware"
	"chatFrame/msggateway"
	"chatFrame/pkg/chatlog"
	"chatFrame/pkg/common/db"
	"chatFrame/pkg/friend"
	"chatFrame/pkg/group"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type TestMode struct {
	Gs        *group.GroupServer
	Userdb    db.UserModelInterface
	Fs        *friend.FriendServer
	Clog      *chatlog.ChatLogServer
	midServer *msggateway.MidServer
}

func (t *TestMode) InitTestRouter(r *gin.Engine) {

	r.Use(middleware.Cors())
	t.Userdb = db.NewUserGorm(db.DBConn)
	t.Gs = group.NewGroupServer()
	t.Fs = friend.NewFriendServer()
	t.Clog = chatlog.NewChatLogServer()
	t.midServer = msggateway.NewMidServer()
	r.POST("createUserAndGroup", t.createUserAndGroup)
	r.POST("joinGroup", t.joinGroup)
	r.POST("kickMember", t.kickMember)
	r.POST("quitGroup", t.quitGroup)
	r.POST("applyFriendReq", t.applyFriendReq)
	r.POST("agreeFriendReq", t.agreeFriendReq)
	r.POST("refuseFriendReq", t.refuseFriendReq)
	r.POST("getRecvFriendReq", t.getRecvFriendReq)
	r.POST("getSendFriendReq", t.getSendFriendReq)
	r.POST("deleteFriend", t.deleteFriend)
	r.POST("setFriendRemark", t.setFriendRemark)
	r.POST("getFriendList", t.getFriendList)

}
func (t *TestMode) createUserAndGroup(c *gin.Context) {

	ginfo := &group.GroupInfo{
		GroupName:    "group test",
		Introduction: "this is a group",
		FaceURL:      "/",
		GroupType:    0,
	}
	m := map[string]interface{}{
		"UserID":   1,
		"NickName": "chd",
		"FaceURL":  "/",
	}
	usermodel := &db.UserModel{
		UserID:     1,
		NickName:   "chd",
		FaceURL:    "/",
		Token:      token.GetToken(m),
		CreateTime: time.Now(),
	}

	if err := t.Userdb.Create(context.Background(), []*db.UserModel{usermodel}); err != nil {
		log.Println(err)
	}
	fmt.Println("Create user")
	if _, err := t.Gs.CreateGroup(ginfo, []uint64{usermodel.UserID}); err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (t *TestMode) joinGroup(c *gin.Context) {
	usermodel := &db.UserModel{
		UserID:     2,
		NickName:   "lvy",
		FaceURL:    "/lvy",
		CreateTime: time.Now(),
	}
	if err := t.Userdb.Create(context.Background(), []*db.UserModel{usermodel}); err != nil {
		log.Println(err)
	}
	groupIDs, err := t.Gs.GroupDB.FindUserJoinedGroupID(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Gs.JoinGroup(groupIDs[0], 1, usermodel.UserID); err != nil {
		log.Fatal(err)
	}
}
func (t *TestMode) kickMember(c *gin.Context) {
	//groupIDs, err := t.Gs.GroupDB.FindUserJoinedGroupID(context.Background(), 1)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//t.Gs.KickMember(groupIDs[0], 2)
	groupIDs, err := t.midServer.GroupServer.GroupDB.FindUserJoinedGroupID(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}
	t.midServer.KickMember(context.Background(), 2, groupIDs[0])
}
func (t *TestMode) quitGroup(c *gin.Context) {
	//groupIDs, err := t.Gs.GroupDB.FindUserJoinedGroupID(context.Background(), 1)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//t.Gs.QuitGroup(groupIDs[0], 2)
	groupIDs, err := t.midServer.GroupServer.GroupDB.FindUserJoinedGroupID(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}
	t.midServer.QuitGroup(context.Background(), 2, groupIDs[0])

}

func (t *TestMode) applyFriendReq(c *gin.Context) {
	if err := t.Fs.ApplyAddFriend(1, 2, "我是丶"); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}

func (t *TestMode) agreeFriendReq(c *gin.Context) {
	if err := t.Fs.AgreeFriendReq(1, 2); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (t *TestMode) refuseFriendReq(c *gin.Context) {
	if err := t.Fs.RefuseFriendReq(1, 2); err != nil {
		log.Fatal(err)

	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}

func (t *TestMode) getRecvFriendReq(c *gin.Context) {
	list, err := t.Fs.GetFriendReqToMe(2, 0, 10)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   list,
	})

}

func (t *TestMode) getSendFriendReq(c *gin.Context) {
	list, err := t.Fs.GetFriendReqFromMe(2, 0, 10)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   list,
	})

}
func (t *TestMode) deleteFriend(c *gin.Context) {
	if err := t.Fs.DeleteFriend(2, 1); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
func (t *TestMode) getFriendList(c *gin.Context) {
	list, err := t.Fs.GetFriendList(2)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   list,
	})
}
func (t *TestMode) setFriendRemark(c *gin.Context) {
	if err := t.Fs.SetFriendRemark(2, 1, "chd"); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
