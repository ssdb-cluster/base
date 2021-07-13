// Copyright 2021 The SSDB-cluster Authors
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
	fn := "tmp/a.wal"
	os.Remove(fn)

	rand.Seed(time.Now().Unix())

	wal := NewWalWriter(fn)
	for i := 0; i< 10; i ++ {
		val := fmt.Sprintf("v%d a", i)
		wal.Write([]byte(val))
	}
	wal.Close()

	reader := NewWalReader(fn)
	for i := 0; i< 10; i ++ {
		val := fmt.Sprintf("v%d a", i)

		bs, err := reader.Next()
		if err != nil {
			t.Fatal(err)
		}
		if bs == nil {
			break
		}

		if string(bs) != val {
			t.Fatal(string(bs) + " != " + val)
		}
	}
	reader.Close()

	// src := []byte("ab\nc")
	// dst := encode(src)
	// src2 := decode(dst[:len(dst)-1])
	// fmt.Println("src", src)
	// fmt.Println("dst", dst)
	// fmt.Println("rev", src2)

}
