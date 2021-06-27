package util

import (
	"os"
	"os/signal"
	"syscall"
)

// wait for SIGNIT|SIGTERM
func WaitForExit() {
	sig_c := make(chan os.Signal, 1)
	signal.Notify(sig_c, syscall.SIGINT, syscall.SIGTERM)
	signal.Ignore(syscall.SIGPIPE)
	<- sig_c
}
