package chatFrame

import (
	"github.com/chdlvy/go-chatFrame/msggateway"
	"github.com/chdlvy/go-chatFrame/pkg/common/config"
	"log"
	"time"
)

func Default() {
	err := config.InitConfig("config.yaml", "/config/")
	if err != nil {
		log.Println("config init failed", err)
	}

}

func RunWsServer(wsPort int, wsSocketTimeout time.Duration, wsMaxConnNum int) error {
	return msggateway.RunWsAndServer(wsPort, wsSocketTimeout, wsMaxConnNum)
}
