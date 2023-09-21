package msggateway

import "time"

const (
	//Time allowed to read and write message every peer
	pongWait = 30 * time.Second
	// max online people of server at the same time
	maxOnlinePeople = 1000
)
