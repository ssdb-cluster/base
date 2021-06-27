package wal

import (
	"testing"
	"os"
	"fmt"
	"base/util"
)

func TestWalReader(t *testing.T){
	if !util.IsDir("tmp") {
		os.MkdirAll("tmp", 0755)
	}

	wal := NewWalReader("tmp/a.wal")
	defer wal.Close()

	for {
		r, err := wal.Next()
		if err != nil {
			fmt.Println("error", err)
			// repair
			// RepairWalFile("tmp/a.wal")
			break
		}
		if r == nil {
			break
		}

		n := 20
		if n > len(r) {
			n = len(r)
		}
		fmt.Println(util.StringEscape(string(r[0:n])))
	}

}
