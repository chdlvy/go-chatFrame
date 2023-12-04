package api

import (
	"context"
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

	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (userRpcServer *UserRpcServer) Logout(context.Context, *userPb.LogoutReq) (*userPb.LogoutResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (userRpcServer *UserRpcServer) Register(ctx context.Context, args *userPb.RegisterReq) (*userPb.RegisterResp, error) {
	hasDBConn()
	userId := utils.GenerateUserID()
	tokenInfo := map[string]interface{}{
		"userName": args.UserName,
		"password": args.Password,
	}
	t := token.GetToken(tokenInfo)
	user := &db.UserModel{UserID: userId, NickName: args.UserName, FaceURL: "no url", Token: t, CreateTime: time.Now()}
	err := userDB.Create(ctx, []*db.UserModel{user})
	if err != nil {

		return &userPb.RegisterResp{}, err
	}
	return &userPb.RegisterResp{AuthToken: t, Code: 1}, nil

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
