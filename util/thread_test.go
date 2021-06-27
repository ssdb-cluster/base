package util

import (
	"testing"
	"log"
	// "time"
)

func TestThread(t *testing.T){
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	task_c := make(chan bool)

	timer := StartTickerThread(200, func(){
		log.Println("ticker")
		task_c <- true
	})

	// go func(){
	// 	for i:=0; i<100; i++ {
	// 		time.Sleep(200 * time.Millisecond)
	// 		if i == 2 {
	// 			log.Println("notify to quit")
	// 			task_c <- false
	// 		}
	// 	}
	// }()

	w := StartWorkerThread(task_c, func(){
		log.Println("work")
	})

	log.Println("")
	Sleep(1)
	log.Println("")
	Sleep(0.1)

	w.Stop()
	log.Println("")
	timer.Stop()
	log.Println("")
	Sleep(2)
}
