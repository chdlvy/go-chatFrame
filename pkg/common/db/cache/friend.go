package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type FriendCache interface {
	GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []uint64, err error)
	AddFriendID(ctx context.Context, ownerUserID string, friendID []string) error
	DelFriendID(ctx context.Context, ownerUserID string, friendID []string) error
	GetFriendCount(ctx context.Context, ownerUserID string) (int64, error)
}

const (
	friendID    = "friend_ids:"
	friendCount = "friend_count:"
)

var _ FriendCache = (*FriendCacheRedis)(nil)

type FriendCacheRedis struct {
	rdb redis.UniversalClient
}

func NewFriendCache(rdb redis.UniversalClient) FriendCache {
	return &FriendCacheRedis{
		rdb: rdb,
	}
}
func (f *FriendCacheRedis) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []uint64, err error) {
	key := friendID + ownerUserID

	res, err := f.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	friendIDs = make([]uint64, len(res))
	for k, v := range res {
		tmp, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		friendIDs[k] = uint64(tmp)
	}
	return friendIDs, nil
}
func (f *FriendCacheRedis) AddFriendID(ctx context.Context, ownerUserID string, friendIDs []string) error {
	key := friendID + ownerUserID
	if err := f.rdb.RPush(ctx, key, friendIDs).Err(); err != nil {
		return err
	}
	countKey := friendCount + ownerUserID
	if _, err := f.rdb.IncrBy(ctx, countKey, int64(len(friendIDs))).Result(); err != nil {
		return err
	}

	return nil

}

func (f *FriendCacheRedis) DelFriendID(ctx context.Context, ownerUserID string, friendIDs []string) error {
	key := friendID + ownerUserID
	count, err := f.rdb.LRem(ctx, key, 0, friendIDs).Result()
	if err != nil {
		return err
	}
	countKey := friendCount + ownerUserID
	if _, err := f.rdb.DecrBy(ctx, countKey, count).Result(); err != nil {
		return err
	}
	return nil

}
func (f *FriendCacheRedis) GetFriendCount(ctx context.Context, ownerUserID string) (int64, error) {
	key := friendCount + ownerUserID
	count, err := f.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(count)
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}
