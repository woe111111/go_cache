package http

import (
	"go_cache/cacheClient"
	"log"
	"net/http"
	"sync"
)

type rebalanceHandler struct {
	*Server
}

var lock sync.RWMutex
var client cacheClient.Client

func (h *rebalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("test")
	if client == nil {
		client = cacheClient.New(h.Addr())
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	lock.Lock()
	go h.rebalance()
}

func (h *rebalanceHandler) rebalance() {
	s := h.NewScanner()
	defer func() {
		s.Close()
		client.RunPipe()
		lock.Unlock()
	}()
	//c := &http.Client{}
	for s.Scan() {
		k := s.Key()
		_, ok := h.ShouldProcess(k)
		if !ok {
			//r, _ := http.NewRequest(http.MethodPut, "http://"+n+":12345/cache/"+k, bytes.NewReader(s.Value()))
			//c.Do(r)
			client.Pipe(&cacheClient.Cmd{"set", k, string(s.Value()), nil})
			h.Del(k)
		}
	}
	return
}

func (s *Server) rebalanceHandler() http.Handler {
	return &rebalanceHandler{s}
}
