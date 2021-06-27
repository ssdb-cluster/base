package wal

import (
	// "os"
)

func RepairWalFile(path string) error {
	// wal := NewWalReader(path)
	// defer wal.Close()

	// pos := wal.Position()
	// for wal.Next() {
	// 	_, err := wal.Item()
	// 	if err != nil {
	// 		os.Truncate(path, pos)
	// 		break
	// 	}
	// 	pos = wal.Position()
	// }

	return nil
}
