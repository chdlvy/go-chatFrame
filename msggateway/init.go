package msggateway

func RunWsAndServer() error {
	if err := StartMidServer(); err != nil {
		return err
	}
	longServer := NewWsServer(Config)
	hubServer := NewHubServer(longServer)
	return hubServer.LongConnServer.Run()
}
