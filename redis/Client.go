// Copyright 2021 The SSDB-cluster Authors
package redis

import (
	"fmt"
	"net"
	"bytes"
)

type Client struct {
	conn *net.TCPConn
}

func NewClient(ip string, port int) *Client {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil
	}

	ret := new(Client)
	ret.conn = conn
	return ret
}

func (cli *Client)Close() {
	cli.conn.Close()
}

func (cli *Client)Do(args ...string) (*Response, error) {
	var err error
	req := NewRequest(args)
	str := req.Encode()
	_, err = cli.conn.Write([]byte(str))
	if err != nil {
		return nil, err
	}

	resp := new(Response)
	var buf bytes.Buffer
	tmp := make([]byte, 32*1024)

	for {
		n := resp.Decode(buf.Bytes())
		if n == -1 {
			recv := string(buf.Bytes())
			if len(recv) > 128 {
				recv = recv[0 : 128]
			}
			err := fmt.Errorf("parse error: %q", recv)
			return nil, err
		} else if n == 0 {
			n, err := cli.conn.Read(tmp)
			if err != nil {
				return nil, err
			}
			buf.Write(tmp[0:n])
		} else {
			return resp, nil
		}
	}
}
