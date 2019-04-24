package cache

import (
	"stathat.com/c/consistent"
)

var shardNum = 32

type conCurrentInMemoryCache struct {
	*consistent.Consistent
	cache map[string]*inMemoryCache
	Stat
}

func (c *conCurrentInMemoryCache) Set(k string, v []byte) error {

	shard, _ := c.Consistent.Get(k)

	c.cache[string(shard)].Set(k, v)

	c.add(k, v) //状态对应更改
	return nil
}

func (c *conCurrentInMemoryCache) Get(k string) ([]byte, error) {

	shard, _ := c.Consistent.Get(k)

	return c.cache[string(shard)].Get(k)
}

func (c *conCurrentInMemoryCache) Del(k string) error {

	shard, _ := c.Consistent.Get(k)

	v, exist := c.cache[string(shard)].c[k]

	if exist {
		delete(c.cache[string(shard)].c, k)
		c.del(k, v.v)
	}
	return nil
}

func (c *conCurrentInMemoryCache) GetStat() Stat {
	return c.Stat
}

func newConCurrentInMemoryCache(ttl int) *conCurrentInMemoryCache {

	circle := consistent.New()
	circle.NumberOfReplicas = shardNum
	var i = 0

	shards := make([]string, shardNum)

	cache := make(map[string]*inMemoryCache, shardNum)

	for {
		shards[i] = string(i)
		i = i + 1

		cache[string(i)] = newInMemoryCache(ttl)

		if i >= shardNum {
			break
		}

	}

	circle.Set(shards)

	return &conCurrentInMemoryCache{circle, cache, Stat{}}

}
