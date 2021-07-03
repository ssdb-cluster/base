package redis

import (
	"base/log"
	"testing"
)

func TestServer(t *testing.T){
	log.SetLevel("debug")

	rport := NewServer("127.0.0.1", 9000)
	defer rport.Close()

	for {
		select {
		case req := <- rport.C:
			resp := new(Response)
			resp.Dst = req.Src
			rport.Send(resp)
		}
	}
}
