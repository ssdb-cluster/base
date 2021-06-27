package wal

import (
	"bytes"
	"strconv"
	"errors"
	"base/util"
)

var err_decode = errors.New("wal corruption")

// 转义 \r\n
func escape_crlf(bs []byte) []byte {
	var buf bytes.Buffer
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
	if s == 0 && e == len(bs) - 1 {
		return bs // no copy
	}
	if s <= e {
		buf.Write(bs[s : e + 1])
	}
	return buf.Bytes()
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
