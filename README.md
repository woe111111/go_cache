# Varys
分布式缓存 提供http,tcp服务  只接受key,value形式的数据

### 外部接口 
```go
type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
	NewScanner() Scanner
}
```
#### Set
- 插入缓存
#### Get 
- 获取缓存数据
#### Del
- 删除数据
#### GetStat
- 获取缓存状态
#### NewScanner
- 集群内部均衡数据
    
### 内存模型

#### conCurrentInMemoryCache
```go
type conCurrentInMemoryCache struct {
	*consistent.Consistent
	cache map[string]*inMemoryCache
	Stat
}
```
通过一致性哈希将key分配到不同基础内存分片上可提高多核并发能力

#### inMemoryCache
```Go
type inMemoryCache struct {
	c     map[string]value
	mutex sync.RWMutex
	Stat
	ttl time.Duration
}
```
通过 map 存储数据 并加锁保证 map 并发安全

### 集群结构
#### 客户端
由于缓存服务是一个需要低延迟的服务 因此服务端没有主节点做负载均衡
通过一致性哈希算法那 key做处理将 将key 传入对应的数据节点到达负载均衡的效果

#### 服务端
服务端内部通过gossip协议保持内部联系，当有集群内部有集群挂了或者有新的机器加入
各自数据节点会自动负载均衡
 
### Use

### HTTP 服务

- put /cache/<key>
- get /cache/<key>
- dellet /cache/<key>
- get /status
- get /cluster
- get /rebalance

### tcp服务
```go
package main

import (
	"Varys/cacheClient"
	"log"
)

func main() {

	res := cacheClient.New("127.0.0.1")


	//cmd := &cacheClient.Cmd{"set","test","1232131",nil}
	//res.Run(cmd)
	//log.Println(cmd.Value)


	res.Pipe(&cacheClient.Cmd{"get","test","",nil})

	log.Println(*res.RunPipe()[0])

}


```

