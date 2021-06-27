package redis

import (
	"testing"
	"fmt"
)

func TestServer(t *testing.T){
	rport := NewServer("127.0.0.1", 9000)
	defer rport.Close()

	for {
		select {
		case msg := <- rport.C:
			fmt.Println(msg.Array())
		}
	}
}
