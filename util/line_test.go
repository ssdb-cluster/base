package util

import (
	"strings"
	"testing"
	"log"
)

func TestLine(t *testing.T){
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	src := "get a\""
	arr := ParseCommandLine(src + "\n")
	log.Println(len(arr), strings.Join(arr, ", "))
}

