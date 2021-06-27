package redis

import (
	"log"
	"testing"
)

func check(t *testing.T, r1, r2 *Response) {
	n := r2.Decode([]byte(r1.Encode() + "sdafsa"))
	if n != len(r1.Encode()) {
		log.Println("error!", n, len(r1.Encode()))
	}
	if r1.Encode() == r2.Encode() {
		log.Println("ok")
	} else {
		t.Error("fail")
		log.Printf("    %q\n", r1.Encode())
		log.Printf("    %q\n", r2.Encode())
	}
}

func TestResponse(t *testing.T){
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	r1 := new(Response)
	r2 := new(Response)

	r1.SetArray([]string{"a", "b"})
	check(t, r1, r2)

	r1.SetInt(123)
	check(t, r1, r2)

	r1.SetNull()
	check(t, r1, r2)

	r1.SetError("message")
	check(t, r1, r2)
}
