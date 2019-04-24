package cacheClient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	server string
}

func (c *httpClient) get(key string) string {
	resp, e := c.Get(c.server + "cache" + key)
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		panic(e)
	}
	return string(b)
}

func (c *httpClient) set(key, value string) {
	req, e := http.NewRequest(http.MethodPut,
		c.server+"cache"+key, strings.NewReader(value))
	if e != nil {
		log.Println(key)
		panic(e)
	}
	resp, e := c.Do(req)
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

func (c *httpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		cmd.Value = c.get(cmd.Key)
		return
	}
	if cmd.Name == "set" {
		c.set(cmd.Key, cmd.Value)
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

func newHTTPClient(server string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://" + server + ":12345/"}
}

func (c *httpClient) Members() ([]string, error) {
	resp, e := c.Get(c.server + "cluster")
	if e != nil {
		return []string{}, e
	}

	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		panic(e)
	}
	str := []string{}
	json.Unmarshal(b, &str)
	return str, nil

}

func (c *httpClient) Pipe(*Cmd) {
	panic("httpClient pipelined run not implement")
}
func (c *httpClient) RunPipe() []*Cmd {
	panic("httpClient pipelined run not implement")
}
