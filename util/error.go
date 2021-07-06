package util

import (
	"io"
	"strings"
)

func IsEOF(err error) bool {
	if err == nil {
		return false
	}
	if err == io.EOF {
		return true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}
	if strings.Contains(err.Error(), "connection reset by peer") {
		return true
	}
	return false
}