package msggateway

import "time"

func RunWsAndServer(wsPort int, wsSocketTimeout time.Duration, wsMaxConnNum int) error {
	if err := StartMidServer(); err != nil {
		return err
	}
	longServer := NewWsServer(wsPort, wsSocketTimeout, wsMaxConnNum)
	hubServer := NewHubServer(longServer)
	return hubServer.LongConnServer.Run()
}
