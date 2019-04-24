package cacheClient

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r      *bufio.Reader
	cmds   []*Cmd
	flag   int
	server string
}

func (c *tcpClient) sendGet(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (c *tcpClient) sendSet(key, value string) {
	klen := len(key)
	vlen := len(value)
	c.Write([]byte(fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)))
}

func (c *tcpClient) sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func readLen(r *bufio.Reader) (int, error) {
	tmp, e := r.ReadString(' ')
	if e != nil {
		return 0, e
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Println(tmp, e)
		return 0, nil
	}
	return l, nil
}

func (c *tcpClient) recvResponse() (string, error) {
	vlen, err := readLen(c.r)
	if err != nil {
		return "", err
	}

	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(c.r, err)
		if e != nil {
			return "", e
		}
		return "", errors.New(string(err))
	}
	value := make([]byte, vlen)
	_, e := io.ReadFull(c.r, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}

func (c *tcpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		c.sendGet(cmd.Key)
		cmd.Value, cmd.Error = c.recvResponse()
		if cmd.Error != nil {
			c.Close()
			c = newTCPClient(c.server)
		}
		return
	}
	if cmd.Name == "set" {
		c.sendSet(cmd.Key, cmd.Value)
		_, cmd.Error = c.recvResponse()
		if cmd.Error != nil {
			c.Close()
			ser, _ := c.Members()
			c = newTCPClient(ser[0])
		}
		return
	}
	if cmd.Name == "del" {
		c.sendDel(cmd.Key)
		_, cmd.Error = c.recvResponse()
		if cmd.Error != nil {
			c.Close()
			ser, _ := c.Members()
			c = newTCPClient(ser[0])
		}
		return
	}

	panic("unknown cmd name " + cmd.Name)
}
func (c *tcpClient) Members() ([]string, error) {
	panic("have not ")
}

func (c *tcpClient) Pipe(cmd *Cmd) {

	if c.flag == 0 {
		c.cmds = make([]*Cmd, 0)
	}
	c.cmds = append(c.cmds, cmd)
	c.flag = 1
}
func (c *tcpClient) RunPipe() []*Cmd {
	if len(c.cmds) == 0 {
		return nil
	}
	for _, cmd := range c.cmds {
		if cmd.Name == "get" {
			c.sendGet(cmd.Key)
		}
		if cmd.Name == "set" {
			c.sendSet(cmd.Key, cmd.Value)
		}
		if cmd.Name == "del" {
			c.sendDel(cmd.Key)
		}
	}
	for _, cmd := range c.cmds {
		cmd.Value, cmd.Error = c.recvResponse()
	}
	c.flag = 0
	return c.cmds

}

func newTCPClient(server string) *tcpClient {
	c, e := net.Dial("tcp", server+":12346")
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)

	cmds := make([]*Cmd, 0)

	return &tcpClient{c, r, cmds, 0, server}
}

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}
