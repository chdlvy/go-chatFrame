package friend

import (
	"github.com/chdlvy/go-chatFrame/pkg/common/db"
	"github.com/chdlvy/go-chatFrame/pkg/common/db/cache"
	"google.golang.org/grpc"
	"log"
)

type friendServer struct {
	friendDatabase db.FriendDatabase
}

func Start() error {
	dbconn, err := db.NewGormDB()
	if err != nil {
		log.Println(err)
	}
	friendDB := db.NewFriendGorm(dbconn)
	friendRequestDB := db.NewFriendRequestGorm(dbconn)
	rdb, err := cache.NewRedis()
	if err != nil {
		log.Fatal(err)
	}
	friendServer := &friendServer{
		friendDatabase: db.NewFriendDatabase(friendDB, friendRequestDB, friendRdb),
	}
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService()
	return friendServer
}
