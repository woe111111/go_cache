package main

import (
	"flag"
	"go_cache/cache"
	"go_cache/cluster"
	"go_cache/http"
	"go_cache/tcp"
	"log"
)

func main() {
	typ := flag.String("type", "concurrent_inmemory", "cache type")
	ttl := flag.Int("ttl", 30, "cache time to live")
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", *typ)
	log.Println("ttl is", *ttl)
	log.Println("node is", *node)
	log.Println("cluster is", *clus)
	c := cache.New(*typ, *ttl)
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	tcpSer := tcp.New(c, n)
	go tcpSer.Listen()

	httpSer := http.New(c, n)

	go httpSer.Listen()
	select {}
	//signal.Listen(tcpSer, httpSer)

}
