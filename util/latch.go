// Copyright 2021 The SSDB-cluster Authors
package util

import (
	// "log"
)

// Latch: level triggered, opaque latch(FlipFlop: edge triggered)
type Latch struct {
	quit_c chan bool
	task_c chan bool
	handler func()
}

func NewLatch() *Latch {
	s := new(Latch)
	s.quit_c = make(chan bool)
	s.task_c = make(chan bool, 1)
	return s
}

func (s *Latch)Close() {
	if s.handler != nil {
		s.task_c <- false
		<- s.quit_c
	}
	close(s.task_c)
	close(s.quit_c)
}

// Bind a signal handler
func (s *Latch)Signal(handler func()) {
	if s.handler != nil {
		s.task_c <- false
		<- s.quit_c
	}
	s.handler = handler

	start_c := make(chan bool)
	go func() {
		start_c <- true

		for {
			t := <- s.task_c
			if t == false {
				break
			}
			s.handler()
		}

		s.quit_c <- true
	}()
	<- start_c
}

// non-blocking
func (s *Latch)Notify() {
	select {
	case s.task_c <- true:
	default:
	}
}
