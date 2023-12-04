package rpc

import "github.com/chdlvy/go-chatFrame/internal/rpc/api"

func Start() {
	api.RunUserRpc()
}
