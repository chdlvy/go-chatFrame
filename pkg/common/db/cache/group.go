package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

const (
	groupMemberIDs = "group_member_ids:"
	groupCount     = "group_count:"
)

type GroupCache interface {
	GetMemberIDs(ctx context.Context, groupID string) (MemberIDs []uint64, err error)
	AddMemberIDs(ctx context.Context, groupID string, MemberIDs []string) error
	DelMemberIDs(ctx context.Context, groupID string, MemberIDs []string) error
	GetGroupCount(ctx context.Context, groupID string) (int64, error)
}
type GroupCacheRedis struct {
	rdb redis.UniversalClient
}

func NewGroupCache(rdb redis.UniversalClient) GroupCache {
	return &GroupCacheRedis{
		rdb: rdb,
	}
}
func (g *GroupCacheRedis) GetMemberIDs(ctx context.Context, groupID string) (MemberIDs []uint64, err error) {
	key := groupMemberIDs + groupID
	res, err := g.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	MemberIDs = make([]uint64, len(res))
	for k, v := range res {
		tmp, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		MemberIDs[k] = uint64(tmp)
	}
	return MemberIDs, nil
}

func (g *GroupCacheRedis) AddMemberIDs(ctx context.Context, groupID string, MemberIDs []string) error {
	key := groupMemberIDs + groupID
	if err := g.rdb.RPush(ctx, key, MemberIDs).Err(); err != nil {
		return err
	}
	countKey := groupCount + groupID
	if _, err := g.rdb.IncrBy(ctx, countKey, int64(len(MemberIDs))).Result(); err != nil {
		return err
	}

	return nil
}

func (g *GroupCacheRedis) DelMemberIDs(ctx context.Context, groupID string, MemberIDs []string) error {
	key := groupMemberIDs + groupID
	count, err := g.rdb.LRem(ctx, key, 0, MemberIDs).Result()
	if err != nil {
		return err
	}
	countKey := groupCount + groupID
	if _, err := g.rdb.DecrBy(ctx, countKey, count).Result(); err != nil {
		return err
	}
	return nil
}

func (g *GroupCacheRedis) GetGroupCount(ctx context.Context, groupID string) (int64, error) {
	key := groupCount + groupID
	count, err := g.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(count)
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}
