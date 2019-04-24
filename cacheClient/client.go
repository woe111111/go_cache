package cacheClient

import (
	"stathat.com/c/consistent"

	"time"
)

type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}
type Client interface {
	Run(cmd *Cmd)
	Pipe(cmd *Cmd)
	RunPipe() []*Cmd
	Members() ([]string, error)
}

type CacheClient struct {
	*consistent.Consistent

	tcp  map[string]Client
	http map[string]Client
	cmds []*Cmd
	flag int
}

func (b *CacheClient) Members() ([]string, error) {
	return b.Consistent.Members(), nil
}

func (b *CacheClient) Run(cmd *Cmd) {
	nodes, err := b.Get(cmd.Key)
	if err != nil {
		panic(err)
	}

	b.tcp[nodes].Run(cmd)

}

func New(server string) Client {

	circle := consistent.New()
	circle.NumberOfReplicas = 256

	tcp := make(map[string]Client)
	http := make(map[string]Client)

	cmds := make([]*Cmd, 0)

	nodes := make([]string, 1)
	nodes[0] = server
	tcp[server] = newTCPClient(server)
	http[server] = newHTTPClient(server)
	circle.Set(nodes)

	go func() {
		for {

			for _, server := range http {
				m, err := server.Members()
				if err != nil {
					continue
				}
				nodes := make([]string, len(m))
				for i, n := range m {

					nodes[i] = n

					_, exist := tcp[n]
					if exist {

					} else {
						tcp[n] = newTCPClient(n)
						http[n] = newHTTPClient(n)
					}
				}
				circle.Set(nodes)
				time.Sleep(time.Second)
			}

		}

	}()

	return &CacheClient{circle, tcp, http, cmds, 0}

}

func (b *CacheClient) Pipe(cmd *Cmd) {
	nodes, err := b.Get(cmd.Key)
	if err != nil {
		panic(err)
	}
	b.tcp[nodes].Pipe(cmd)

	if b.flag == 0 {
		b.cmds = make([]*Cmd, 0)
	}
	b.cmds = append(b.cmds, cmd)
	b.flag = 1

}

func (b *CacheClient) RunPipe() []*Cmd {

	for _, ser := range b.tcp {
		ser.RunPipe()
	}
	b.flag = 0
	return b.cmds
}
