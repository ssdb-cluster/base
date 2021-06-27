package redis

import (
	"log"
	"testing"
)

func TestClient(t *testing.T){
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	client := NewClient("127.0.0.1", 8888)
	defer client.Close()

	var resp *Response
	var err error

	resp, err = client.Do("set", "a", "你好")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if resp.IsError() {
		t.Error(resp.ErrorCode(), resp.ErrorMessage())
	}

	resp, err = client.Do("get", "a")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if resp.IsNull() {
		log.Println("not found")
	}
	if resp.IsString() {
		log.Println(resp.String())
	}

	resp, err = client.Do("scan", "Z", "", "3")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if resp.IsArray() {
		log.Println(resp.Pairs())
	}

}
