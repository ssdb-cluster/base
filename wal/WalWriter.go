// Copyright 2021 The SSDB-cluster Authors
package wal

import (
	"os"
	"fmt"
	"bufio"
)

var WalWriterBufferSize int = 4 * 1024 * 1024

type WalWriter struct{
	fp *os.File
	path string
	size int64
	writer *bufio.Writer
}

// create if not exists
func NewWalWriter(path string) *WalWriter {
	var size int64
	fp, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	size, err = fp.Seek(0, os.SEEK_END)
	if err != nil {
		fmt.Println(err)
		fp.Close()
		return nil
	}

	ret := new(WalWriter)
	ret.fp = fp
	ret.path = path
	ret.size = size
	ret.writer = bufio.NewWriterSize(ret.fp, WalWriterBufferSize)

	return ret
}

func (wal *WalWriter)Close() {
	if err := wal.Fsync(); err != nil {
		panic(err)
	}
	if err := wal.fp.Close(); err != nil {
		panic(err)
	}
}

func (wal *WalWriter)Path() string {
	return wal.path
}

func (wal *WalWriter)Size() int64 {
	return wal.size
}

// Truncate and fsync
func (wal *WalWriter)Truncate(size int64) error {
	if wal.writer != nil {
		if err := wal.writer.Flush(); err != nil {
			return err
		}
	}
	if err := wal.fp.Truncate(size); err != nil {
		return err
	}
	if err := wal.fp.Sync(); err != nil {
		return err
	}
	wal.size = size
	_, err := wal.fp.Seek(size, os.SEEK_SET)
	return err
}

func (wal *WalWriter)Write(bs []byte) (nn int64, err error) {
	buf := encode(bs)
	var size int
	if wal.writer != nil {
		size, err = wal.writer.Write(buf)
	} else {
		size, err = wal.fp.Write(buf)
	}

	wal.size += int64(size)
	return int64(size), err
}

func (wal *WalWriter)Fsync() error {
	if wal.writer != nil {
		if err := wal.writer.Flush(); err != nil {
			return err
		}
	}
	return wal.fp.Sync()
}
