package msggateway

import (
	"fmt"
	"sync"
)

type UserMap struct {
	m sync.Map
}

func newUserMap() *UserMap {
	return &UserMap{}
}
func (u *UserMap) GetAll(key string) ([]*Client, bool) {
	allClients, ok := u.m.Load(key)
	if ok {
		return allClients.([]*Client), ok
	}
	return nil, ok
}
func (u *UserMap) Get(key string, platformID int) ([]*Client, bool) {
	fmt.Println(key)

	allClients, userExisted := u.m.Load(key)
	fmt.Println(userExisted)
	if userExisted {
		var clients []*Client
		for _, client := range allClients.([]*Client) {
			if client.PlatformID == platformID {
				clients = append(clients, client)
			}
		}
		return clients, userExisted
	}
	return nil, userExisted
}
func (u *UserMap) Set(key string, v *Client) {
	allClients, existed := u.m.Load(key)
	if existed {
		oldClients := allClients.([]*Client)
		oldClients = append(oldClients, v)
		u.m.Store(key, oldClients)
	} else {
		var clients []*Client
		clients = append(clients, v)
		u.m.Store(key, clients)
	}
}
func (u *UserMap) DeleteAll(key string) {
	u.m.Delete(key)
}
