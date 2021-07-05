package util

import (
	"time"
	"sync"
	// "log"
)

type Thread interface {
	Stop()
}

// two ways to stop the thread:
// 1. task_c <- false
// 2. Stop()
func StartWorkerThread(task_c chan bool, f func()) Thread {
	// log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	w := new(WorkerThread)
	w.task_c = task_c
	w.work_f = f
	w.start()
	return w
}

// call Stop() to stop the thread
func StartTickerThread(tick_ms int, f func()) *TickerThread {
	w := new(TickerThread)
	w.tick_ms = tick_ms
	w.work_f = f
	w.start()
	return w
}

type WorkerThread struct {
	sync.Mutex
	running bool
	done_c chan bool
	task_c chan bool
	work_f func()
}

func (w *WorkerThread)start() {
	w.done_c = make(chan bool, 1)
	w.running = true

	go func() {
		defer func(){
			for len(w.task_c) > 0 {
				<- w.task_c
			}
			w.done_c <- true

			w.Lock()
			w.running = false
			w.Unlock()
		}()

		for {
			t := <- w.task_c
			if t == false {
				return
			}
			w.work_f()
		}
	}()
}

func (w *WorkerThread)Stop() {
	w.Lock()
	defer w.Unlock()

	if w.running == false {
		return
	}
	w.task_c <- false
	<- w.done_c
}

type TickerThread struct {
	ticker *time.Ticker
	stop_c chan bool
	work_f func()
	tick_ms int
}

func (w *TickerThread)start() {
	w.stop_c = make(chan bool)
	w.ticker = time.NewTicker(time.Duration(w.tick_ms) * time.Millisecond)

	go func(){
		for {
			select{
			case <- w.stop_c:
				return
			case <- w.ticker.C:
				w.work_f()
			}
		}
	}()
}

func (w *TickerThread)Stop() {
	w.ticker.Stop()
	w.stop_c <- false
}
