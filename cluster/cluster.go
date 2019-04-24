package cluster

import (
	"github.com/hashicorp/memberlist"
	"io/ioutil"
	"log"
	"net/http"
	"stathat.com/c/consistent"
	"time"
)

type Node interface {
	ShouldProcess(key string) (string, bool)
	Members() []string
	Addr() string
	Leave()
}

type node struct {
	*consistent.Consistent
	l *memberlist.Memberlist

	addr string
}

func (n *node) Addr() string {
	return n.addr
}

var menbers = []*memberlist.Node{}

func (n *node) Members() []string {
	return n.Consistent.Members()
}

func New(addr, cluster string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}
	circle := consistent.New()
	circle.NumberOfReplicas = 256

	client := http.Client{}
	go func() {
		for {

			m := l.Members()

			if len(m) != len(menbers) {
				client.Get("http://" + addr + ":12345/rebalance")
			}

			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)

			menbers = m
		}
	}()
	return &node{circle, l, addr}, nil
}

func (n *node) Leave() {
	n.l.Leave(1)
	log.Println(n.Members())
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}
