package cache

import "log"

func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemoryCache(ttl)
	}
	//if typ == "rocksdb" {
	//	c = newRocksdbCache(ttl)
	//}

	if typ == "concurrent_inmemory" {
		c = newConCurrentInMemoryCache(ttl)
	}

	if c == nil {
		panic("unknown cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
