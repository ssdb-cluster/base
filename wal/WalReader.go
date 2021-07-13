// Copyright 2021 The SSDB-cluster Authors
package wal

import (
	"io"
	"os"
	"fmt"
	"errors"
	"bufio"
	"bytes"
)

// supports both random read and sequential read
type WalReader struct{
	fp *os.File
	reader *bufio.Reader
	path string
	fpos int64
	size int64
}

func NewWalReader(path string) *WalReader {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	var size int64
	size, err = fp.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	_, err = fp.Seek(0, os.SEEK_SET)
	if err != nil {
		panic(err)
	}

	ret := new(WalReader)
	ret.fp = fp
	ret.path = path
	ret.size = size
	ret.reader = bufio.NewReader(ret.fp)

	return ret
}

func (wal *WalReader)Close() {
	if err := wal.fp.Close(); err != nil {
		panic(err)
	}
}

func (wal *WalReader)Path() string {
	return wal.path
}

func (wal *WalReader)Size() int64 {
	return wal.size
}

// current read position
func (wal *WalReader)Position() int64 {
	return wal.fpos
}

func (wal *WalReader)Fseek(pos int64) {
	_, err := wal.fp.Seek(pos, os.SEEK_SET)
	if err != nil {
			panic(err)
	}
	wal.reader.Reset(wal.fp)
	wal.fpos = pos
}

func (wal *WalReader)First() ([]byte, error) {
	wal.Fseek(0)
	return wal.Next()
}

func (wal *WalReader)Last() ([]byte, error) {
	wal.Fseek(0)
	var ret []byte
	for {
		bs, err := wal.Next()
		if err != nil {
			return nil, err
		}
		if bs == nil {
			break
		}
		ret = bs
	}
	return ret, nil
}
// read exactly next record, place read position at just after the record
func (wal *WalReader)Next() ([]byte, error) {
	bs, err := wal.reader.ReadBytes('\n')
	wal.fpos += int64(len(bs))
	if err != nil {
		if err == io.EOF {
			if len(bs) > 0 {
				return nil, errors.New("wal corruption")
			}
			return nil, nil
		}
		return nil, err
	}

	bs, err = decode(bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// cause read position undefined
func (wal *WalReader)ReadAt(pos int64) ([]byte, error) {
	start := pos
	var bs []byte
	tmp := make([]byte, 4096)
	for {
		// ReadAt() uses pread(), thread safe
		n, err := wal.fp.ReadAt(tmp, start)
		if n == 0 {
			return nil, fmt.Errorf("read at pos: %d, %v", start, err)
		}
		p := bytes.IndexByte(tmp[0:n], '\n')
		if p != -1 {
			if len(bs) == 0 {
				bs = tmp[0:p+1]
			} else {
				bs = append(bs, tmp[0:p+1]...)
			}
			break
		}
		start += int64(n)
		bs = append(bs, tmp[0:n]...)
	}

	ret, err := decode(bs)
	if err != nil {
		return bs, fmt.Errorf("wal record decode error at pos: %d, %v", pos, err)
	}
	return ret, nil
}
