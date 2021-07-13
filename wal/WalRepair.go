// Copyright 2021 The SSDB-cluster Authors
package wal

import (
	"os"
)

func RepairWalFile(path string) error {
	r := NewWalReader(path)
	defer r.Close()

	for {
		pos := r.Position()
		_, err := r.Next()
		if err != nil {
			if err := os.Truncate(path, pos); err != nil {
				panic(err)
			}
			break
		}
	}

	return nil
}
