package msggateway

func RunWsAndServer() error {
	longServer := NewWsServer(Config)
	hubServer := NewHubServer(longServer)
	return hubServer.LongConnServer.Run()
}
