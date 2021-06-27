package util

import (
	"testing"
	"log"
	"time"
)

func TestLatch(t *testing.T){
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	s := NewLatch()
	s.Notify()
	s.Notify()
	s.Notify()

	if len(s.task_c) != 1 {
		t.Fatal("")
	}

	count := 0
	s.Signal(func(){
		count ++
		log.Println("fire")
	})

	// wait consumer
	time.Sleep(1 * time.Millisecond)

	s.Notify()
	s.Close()
	// wait consumer
	// time.Sleep(1 * time.Millisecond)

	if count != 2 {
		t.Fatal("")
	}
}
