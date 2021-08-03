// Copyright 2021 The SSDB-cluster Authors
package redis

import (
	"bytes"
	"sync"
	"net"
	"errors"
	"base/log"
)

// 如果客户端不是请求响应模式, 而是 pipeline 模式, 可能会遇到这样的情况:
// 第1个请求正确地发往 raft, 第2个请求时 raft 返回 not_leader 错误. 这样的话, 第2个响应会先返回, 从而出错.
// 所以, 通过 recv_wait_c 保证前一个请求的响应发送之后, 才能接收下一个请求
type serv_link_t struct {
	sync.Mutex
	closed bool

	id int
	isRedis bool
	conn net.Conn

	recv_tmp []byte
	recv_buf bytes.Buffer
	resp_c chan *Response
	recv_wait_c chan bool
}

func new_serv_link(id int, conn *net.TCPConn) *serv_link_t {
	conn.SetNoDelay(true)

	client := new(serv_link_t)
	client.id = id
	client.conn = conn
	client.resp_c = make(chan *Response, 1024)
	client.recv_wait_c = make(chan bool)
	client.recv_tmp = make([]byte, 128*1024)

	// 让 recv 能开始执行第一次
	go func() {
		client.recv_wait_c <- true
	}()

	return client
}

func (l *serv_link_t)close() {
	l.Lock()
	defer l.Unlock()

	if l.closed {
		return
	}
	l.closed = true

	l.conn.Close()
	close(l.resp_c)
	close(l.recv_wait_c)
}

func (client *serv_link_t)send(resp *Response) error {
	var data string
	if client.isRedis {
		data = resp.Encode()
	} else {
		data = resp.EncodeSSDB()
	}

	log.Trace("   send > %d %s", resp.Dst, data)
	bs := []byte(data)
	for len(bs) > 0 {
		nn, err := client.conn.Write(bs)
		if err != nil {
			return err
		}
		bs = bs[nn : ]
	}

	client.recv_wait_c <- true
	return nil
}

func (client *serv_link_t)recv() (*Request, error) {
	w := <- client.recv_wait_c
	if w == false {
		return nil, nil
	}

	msg := new(Request)
	for {
		for client.recv_buf.Len() > 0 {
			n := msg.Decode(client.recv_buf.Bytes())
			if n == -1 {
				return nil, errors.New("parse error")
			} else if (n == 0){
				break
			}
			str := string(client.recv_buf.Bytes()[0 : n])
			client.recv_buf.Next(n)

			if len(str) > 100 {
				log.Trace("recv > %s...", str[0:100])
			} else {
				log.Trace("recv > %s", str)
			}

			client.isRedis = msg.IsRedis
			msg.Src = client.id
			return msg, nil
		}

		n, err := client.conn.Read(client.recv_tmp)
		if err != nil {
			return nil, err
		}
		client.recv_buf.Write(client.recv_tmp[0:n])
	}
}
