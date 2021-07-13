// Copyright 2021 The SSDB-cluster Authors
package wal

import (
	"bytes"
	"strconv"
	"errors"
	"base/util"
)

// 转义 \r\n
func escape_crlf(bs []byte, buf *bytes.Buffer) {
	var s int = 0
	var e int = -1
	var c byte
	for e, c = range bs {
		var d string
		switch c {
		case '\\':
			d = "\\\\"
		case '\r':
			d = "\\r"
		case '\n':
			d = "\\n"
		default:
			continue
		}
		buf.Write(bs[s : e])
		buf.WriteString(d)
		s = e + 1
	}
	if s <= e {
		buf.Write(bs[s : e + 1])
	}
}

func hex32(n uint32, buf []byte) {
	const board = "0123456789abcdef"
	shift := uint32(32)
	for i := 0; i < 8; i ++ {
		shift -= 4
		s := ((n >> shift) & 0xf)
		b := board[s]
		buf[i] = b
	}
}

func encode(bs []byte) []byte {
	var buf *bytes.Buffer
	// 根据 benchmark, 这样做可以优化:
	// * make 传常量大小
	// * 减少 if 分支数量
	if len(bs) < 200 {
		buf = bytes.NewBuffer(make([]byte, 8, 256))
	} else if len(bs) < 400 {
		buf = bytes.NewBuffer(make([]byte, 8, 512))
	} else if len(bs) < 800 {
		buf = bytes.NewBuffer(make([]byte, 8, 1024))
	} else {
		nn := len(bs) + len(bs) / 8
		buf = bytes.NewBuffer(make([]byte, 8, nn))
	}

	buf.WriteByte(' ')
	escape_crlf(bs, buf)
	buf.WriteByte('\n')

	ret := buf.Bytes()
	sum := util.Crc32(ret[9 : len(ret) - 1])
	hex32(sum, ret)
	return ret

	// buf := bytes.NewBuffer(make([]byte, 0, nn))
	// escape_crlf(bs, buf)
	// sum := util.Crc32(buf.Bytes())
	// s := fmt.Sprintf("%08x %s\n", sum, buf.Bytes())
	// return []byte(s)
}

func decode(bs []byte) ([]byte, error) {
	if len(bs) < 9 {
		return nil, errors.New("incomplete wal item")
	}
	ret := bs[9 : len(bs)-1]
	sum0 := util.Crc32([]byte(ret))
	sum1, err := strconv.ParseUint(string(bs[0:8]), 16, 32)
	if err != nil {
		return nil, errors.New("checksum parse failed: " + err.Error())
	}
	if sum0 != uint32(sum1) {
		return nil, errors.New("checksum mismatch")
	}

	ret = util.BytesUnescape(ret)
	return ret, nil
}
