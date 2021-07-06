package redis

import (
	"base/util"
	"net"
	"fmt"
	"sync"
	"base/log"
)

/*
TODO:
如果客户端不是请求响应模式, 而是 pipeline 模式, 可能会遇到这样的情况:
第1个请求正确地发往 raft, 第2个请求时 raft 返回 not_leader 错误.
这样的话, 第2个响应会先返回, 从而出错.
*/
type Server struct {
	sync.Mutex
	wait sync.WaitGroup

	recv_c chan *Request

	lastClientId int
	conn *net.TCPListener
	clients map[int]*serv_link_t
}

func NewServer(ip string, port int) *Server {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	tp := new(Server)
	tp.lastClientId = 0
	tp.conn = conn
	tp.recv_c = make(chan *Request, 1024)
	tp.clients = make(map[int]*serv_link_t)

	tp.wait.Add(1)
	go tp.accept_thread()

	log.Info("redis server listen on tcp://%s:%d", ip, port)
	return tp
}

func (tp *Server)Close(){
	tp.Lock()
	for _, client := range tp.clients {
		client.close()
	}
	tp.Unlock()
	tp.conn.Close()
	tp.wait.Wait()

	close(tp.recv_c)
	log.Info("redis server closed")
}

func (tp *Server)C() chan *Request {
	return tp.recv_c
}

func (tp *Server)Send(resp *Response) {
	tp.Lock()
	defer tp.Unlock()

	dst := resp.Dst
	client := tp.clients[dst]
	if client == nil {
		log.Trace("connection not found: %d", dst)
		return
	}

	client.resp_c <- resp
}

func (tp *Server)accept_thread() {
	defer tp.wait.Done()
	for {
		tp.lastClientId ++

		conn, err := tp.conn.AcceptTCP()
		if err != nil {
			if !util.IsEOF(err) {
				log.Error("%v", err)
			}
			return
		}

		client := new_serv_link(tp.lastClientId, conn)
		log.Debug("accept connection %d %s", client.id, conn.RemoteAddr())

		tp.Lock()
		tp.clients[client.id] = client
		tp.Unlock()

		tp.wait.Add(2)
		go tp.recv_thread(client)
		go tp.send_thread(client)
	}
}

func (tp *Server)send_thread(client *serv_link_t) {
	defer tp.wait.Done()
	for {
		resp := <- client.resp_c
		if resp == nil {
			break
		}
		if err := client.send(resp); err != nil {
			client.close()
			if !util.IsEOF(err) {
				log.Error("send error: %v", err)
			}
			return
		}
	}
}

func (tp *Server)recv_thread(client *serv_link_t) {
	defer tp.wait.Done()
	defer func() {
		tp.Lock()
		if tp.clients[client.id] != nil {
			log.Debug("close connection %d %s", client.id, client.conn.RemoteAddr())
			delete(tp.clients, client.id)
		}
		tp.Unlock()
		client.close()
	}()

	for {
		req, err := client.recv()
		if err != nil {
			if !util.IsEOF(err) {
				log.Error("recv error: %v", err)
			}
			return
		}
		if req == nil {
			return
		}
		tp.recv_c <- req
	}
}
