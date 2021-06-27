package util

import (
	"hash/crc32"
)

func Crc32(s []byte) uint32 {
	return crc32.ChecksumIEEE(s)
}
