package msggateway

import (
	"strconv"
	"sync"
)

// store all clients of online
type UserMap struct {
	m sync.Map
}

func newUserMap() *UserMap {
	return &UserMap{}
}
func (u *UserMap) Get(key string) (*Client, bool) {

	client, userExisted := u.m.Load(key)
	if userExisted {
		return client.(*Client), userExisted
	}
	return nil, userExisted
}
func (u *UserMap) GetManyUsers(key []string, clients []*Client) bool {
	for k, v := range key {
		client, userExisted := u.m.Load(v)
		if !userExisted {
			return userExisted
		}
		clients[k] = client.(*Client)
	}

	return true
}
func (u *UserMap) Set(key string, v *Client) {
	u.m.Store(key, v)
}
func (u *UserMap) Delete(key string) {
	u.m.Delete(key)
}
func (u *UserMap) GetAllClient(clients map[uint64]*Client) {
	u.m.Range(func(key, value any) bool {
		userId, _ := strconv.Atoi(key.(string))
		clients[uint64(userId)] = value.(*Client)
		return true
	})
}
