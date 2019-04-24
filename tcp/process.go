package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

type result struct {
	v []byte
	e error
}

func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		v, e := s.Get(k)
		c <- &result{v, e}
	}()
}

func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Set(k, v)}
	}()
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Del(k)}
	}()
}

func reply(conn net.Conn, resultCh chan chan *result) {
	defer conn.Close()
	for {
		c, open := <-resultCh
		if !open {
			return
		}
		r := <-c
		e := sendResponse(r.v, r.e, conn)
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
	}
}

func (s *Server) process(conn net.Conn) {
	r := bufio.NewReader(conn)
	resultCh := make(chan chan *result, 5000)
	// 这里使用 channel  channel *result 而不是channel
	// 因为 若只是channel *result  携程结束任务后往result内写数据 不能保证按照来时的顺序

	//  channel channel *result  来时的访问顺序写入
	// 返回时
	// for{
	//      res_channel <- result   此时顺序取出一个
	//      res  <- res_channel    程序会再次等待 channel *result对象填充  达到顺序返回效果  达到同步效果
	/// }
	defer close(resultCh)

	// 在一个链接内 输入 输出 通过 管道进行解耦
	go reply(conn, resultCh)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}
		if op == 'S' {
			s.set(resultCh, r)
		} else if op == 'G' {
			s.get(resultCh, r)
		} else if op == 'D' {
			s.del(resultCh, r)
		} else {
			log.Println("close connection due to invalid operation:", op)
			return
		}
	}
}
