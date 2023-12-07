package api

import (
	"context"
	"errors"
	token "github.com/chdlvy/go-chatFrame/jwt"
	"github.com/chdlvy/go-chatFrame/pkg/common/config"
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/chdlvy/go-chatFrame/utils"
	userPb "github.com/chdlvy/protocol/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
	"net"
	"time"
)

const (
	SuccessReplyCode = 1
)

type UserRpcServer struct {
	userPb.UnimplementedUserServiceServer
}

func RunUserRpc() {

	listen, _ := net.Listen("tcp", config.Config.Rpc.UserAddress)
	grpcServer := grpc.NewServer()
	userPb.RegisterUserServiceServer(grpcServer, &UserRpcServer{})
	err := grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("run user rpc failed", err)
		return
	}
}

var gdb *gorm.DB
var userDB db.UserModelInterface

func (userRpcServer *UserRpcServer) Login(ctx context.Context, loginReq *userPb.LoginReq) (*userPb.LoginResp, error) {
	hasDBConn()
	if err := userDB.CheckUserByPhone(ctx, loginReq.Phone, loginReq.Password); err != nil {
		return &userPb.LoginResp{}, errors.New("no this phone or password error!")
	}
	tokenInfo := map[string]interface{}{
		"userName": loginReq.Phone,
		"password": loginReq.Password,
	}
	t := token.GetToken(tokenInfo)

	//redis操作.....

	return &userPb.LoginResp{AuthToken: t, Code: SuccessReplyCode}, nil
}
func (userRpcServer *UserRpcServer) Logout(context.Context, *userPb.LogoutReq) (*userPb.LogoutResp, error) {
	//删除redis中的token即可
	return &userPb.LogoutResp{Code: SuccessReplyCode}, nil
}
func (userRpcServer *UserRpcServer) Register(ctx context.Context, args *userPb.RegisterReq) (*userPb.RegisterResp, error) {
	hasDBConn()
	userId := utils.GenerateUserID()
	tokenInfo := map[string]interface{}{
		"userName": args.Phone,
		"password": args.Password,
	}
	t := token.GetToken(tokenInfo)
	user := &db.UserModel{UserID: userId, NickName: args.UserName, FaceURL: "no url", Token: t, Phone: args.Phone, CreateTime: time.Now()}
	err := userDB.Create(ctx, []*db.UserModel{user})
	if err != nil {

		return &userPb.RegisterResp{}, err
	}
	return &userPb.RegisterResp{AuthToken: t, Code: SuccessReplyCode}, nil

}
func (userRpcServer *UserRpcServer) CheckAuth(context.Context, *userPb.CheckAuthReq) (*userPb.CheckAuthResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAuth not implemented")
}

func hasDBConn() {
	if gdb == nil {
		gdb, err := db.NewGormDB()
		if err != nil {
			log.Fatal(err)
		}
		userDB = db.NewUserGorm(gdb)
	}
	if userDB == nil {
		userDB = db.NewUserGorm(gdb)
	}
}
