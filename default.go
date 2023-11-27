package chatFrame

import (
	"github.com/chdlvy/go-chatFrame/internal/msggateway"
	"time"
)

func RunWsServer(wsPort int, wsSocketTimeout time.Duration, wsMaxConnNum int) error {
	return msggateway.RunWsAndServer(wsPort, wsSocketTimeout, wsMaxConnNum)
}
