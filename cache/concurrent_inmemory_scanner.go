package cache

type conCurrent_InMemoryScanner struct {
	pair
	pairCh  chan *pair
	closeCh chan struct{}
}

func (s *conCurrent_InMemoryScanner) Close() {
	close(s.closeCh)
}

func (s *conCurrent_InMemoryScanner) Scan() bool {
	p, ok := <-s.pairCh
	if ok {
		s.k, s.v = p.k, p.v
	}
	return ok
}

func (s *conCurrent_InMemoryScanner) Key() string {
	return s.k
}

func (s *conCurrent_InMemoryScanner) Value() []byte {
	return s.v
}

func (c *conCurrentInMemoryCache) NewScanner() Scanner {
	pairCh := make(chan *pair)
	closeCh := make(chan struct{})
	go func() {
		defer close(pairCh) //
		for _, ca := range c.cache {
			ca.mutex.RLock()
			for k, v := range ca.c {
				ca.mutex.RUnlock()
				select {
				case <-closeCh:
					return
				case pairCh <- &pair{k, v.v}:
				}
				ca.mutex.RLock()
			}
			ca.mutex.RUnlock()
		}

	}()
	//该匿名函数用作并发控制
	//因为在select时程序可能等待通道数据 因此需要先解锁 获取到通道数据后在加锁  可解决并发访问问题
	//通过 两个channel相互关闭从而在两个协程之间通讯

	//协程关闭逻辑
	//当缓存读完后 pairCh 内无数据
	// Scan 返回数据 flase
	// 执行 s.close closeCh 关闭
	// closeCh 可读
	// pairCh 关闭

	return &conCurrent_InMemoryScanner{pair{}, pairCh, closeCh}
}
