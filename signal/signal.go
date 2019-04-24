package signal

import (
	"fmt"
	"go_cache/http"
	"go_cache/tcp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Listen(tcpServer *tcp.Server, httpServer *http.Server) {
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)
	tcpServer.Node.Leave()
	time.Sleep(time.Second * 3)
	// the back scan while do it
	//client := htp.Client{}
	//client.Get("http://127.0.0.1:12345/rebalance")

	for {
		if tcpServer.GetStat().Count == 0 {
			break
		}
		time.Sleep(time.Second)
	}

}
