package wal

import (
	"testing"
	"os"
	"fmt"
	"time"
	"math/rand"
	"base/util"
)

func TestWalWriter(t *testing.T){
	if !util.IsDir("tmp") {
		os.MkdirAll("tmp", 0755)
	}

	rand.Seed(time.Now().Unix())

	wal := NewWalWriter("tmp/a.wal")
	for i := 0; i< 10; i ++ {
		val := fmt.Sprintf("v%d", i)
		wal.Write([]byte(val))
	}
	wal.Close()

	// src := []byte("ab\nc")
	// dst := encode(src)
	// src2 := decode(dst[:len(dst)-1])
	// fmt.Println("src", src)
	// fmt.Println("dst", dst)
	// fmt.Println("rev", src2)

}
