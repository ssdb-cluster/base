package redis

import (
	"io"
	"net"
	"fmt"
	"sync"
	"bytes"
	"strings"
	"base/log"
	"base/util"
)

/*
TODO: 请求响应模式, 一个连接如果有一个请求在处理时, 则不再解析报文, 等响应后再解析下一个报文.
*/
type Server struct {
	sync.Mutex
	C chan *Request

	lastClientId int
	conn *net.TCPListener
	clients map[int]*client_t
	resp_c chan *Response
}

type client_t struct {
	id int
	conn net.Conn
	isRedis bool
}

func NewServer(ip string, port int) *Server {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	tp := new(Server)
	tp.C = make(chan *Request, 1024)
	tp.lastClientId = 0
	tp.conn = conn
	tp.clients = make(map[int]*client_t)
	tp.resp_c = make(chan *Response, 1024)

	// TODO: create goroutine for each connection?
	go func() {
		for {
			resp := <- tp.resp_c
			if resp == nil {
				break
			}
			tp.doSend(resp)
		}
	}()

	log.Info("redis server listen on tcp://%s:%d", ip, port)
	tp.start()
	return tp
}

func (tp *Server)Close(){
	tp.Lock()
	for _, client := range tp.clients {
		client.conn.Close()
	}
	tp.Unlock()
	tp.conn.Close()
	close(tp.C)
	close(tp.resp_c)
	log.Info("redis server closed")
}

func (tp *Server)start() {
	go func(){
		for {
			tp.lastClientId ++

			conn, err := tp.conn.AcceptTCP()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Error("%v", err)
				}
				return
			}
			conn.SetNoDelay(true)
			conn.SetWriteBuffer(1024 * 1024)

			client := new(client_t)
			client.id = tp.lastClientId
			client.conn = conn
			tp.Lock()
			tp.clients[client.id] = client
			tp.Unlock()

			log.Debug("accept connection %d %s", client.id, conn.RemoteAddr().String())
			go tp.receiveClient(client)
		}
	}()
}

func (tp *Server)receiveClient(client *client_t) {
	defer func() {
		log.Debug("close connection %d %s", client.id, client.conn.RemoteAddr().String())
		tp.Lock()
		delete(tp.clients, client.id)
		tp.Unlock()
		client.conn.Close()
	}()

	var buf bytes.Buffer
	var msg *Request
	msg = new(Request)
	tmp := make([]byte, 64*1024)

	for {
		for {
			n := msg.Decode(buf.Bytes())
			if n == -1 {
				recv := string(buf.Bytes())
				if len(recv) > 128 {
					recv = recv[0 : 128]
				}
				log.Warn("%v parse error: %q", client.conn.RemoteAddr(), recv)

				resp := new(Response)
				resp.Dst = client.id
				resp.SetError("parse error")
				client.conn.Write([]byte(resp.Encode()))
				return
			} else if (n == 0){
				break
			}
			log.Trace("receive < %d %s", client.id, util.StringEscape(string(tmp[0:n])))
			buf.Next(n)

			msg.Src = client.id
			client.isRedis = msg.IsRedis

			tp.C <- msg
			msg = new(Request)
		}

		n, err := client.conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				return
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			if strings.Contains(err.Error(), "connection reset by peer") {
				return
			}
			log.Error("%v", err)
			return
		}
		buf.Write(tmp[0:n])
	}
}

func (tp *Server)Send(resp *Response) {
	tp.resp_c <- resp
}

func (tp *Server)doSend(resp *Response) {
	tp.Lock()
	defer tp.Unlock()

	dst := resp.Dst
	client := tp.clients[dst]
	if client == nil {
		log.Trace("connection not found: %d", dst)
		return
	}

	var data string
	if client.isRedis {
		data = resp.Encode()
	} else {
		data = resp.EncodeSSDB()
	}

	log.Trace("   send > %d %s\n", dst, util.StringEscape(data))
	client.conn.Write([]byte(data))
}
