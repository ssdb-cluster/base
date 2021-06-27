package redis

import (
	"bytes"
	"strings"
	"strconv"
)

type Reply struct {
	type_ int
	err string // when ERROR
	num int    // when INT
	val string // when INT, STRING, ERROR, ARRAY
}

func (reply *Reply)Decode(bs []byte) int {
	idx, val := read_line(bs)
	if idx == 0 {
		return 0
	}
	val = string(val[1 : ])

	switch bs[0] {
	case '+':
		if val != "OK" {
			return -1
		}
		reply.type_ = TypeOK
	case '-':
		reply.type_ = TypeError
		ps := strings.SplitN(val, " ", 2)
		reply.err = ps[0]
		if len(ps) == 2 {
			reply.val = ps[1]
		}
	case ':':
		reply.type_ = TypeInt
		reply.num, _ = strconv.Atoi(val)
		reply.val = val
	case '$':
		size, _ := strconv.Atoi(val)
		if size == -1 {
			reply.type_ = TypeNull
		} else {
			reply.type_ = TypeString
			parsed, data := read_bytes(bs[idx: ], size)
			if parsed <= 0 {
				return 0
			}
			reply.val = string(data)
			idx += parsed
		}
	case '*':
		reply.type_ = TypeArray
		reply.num, _ = strconv.Atoi(val)
	default:
		return -1
	}

	return idx
}

// returned val does not include \n or \r\n
func read_line(bs []byte) (parsed int, val string) {
	idx := bytes.IndexByte(bs, '\n')
	if idx == -1 {
		return 0, ""
	}
	p := bs[0 : idx]
	if len(p) > 0 && p[len(p)-1] == '\r' {
		p = p[0 : len(p)-1]
	}
	idx += 1
	return idx, string(p)
}

// bytes are ended with \n or \r\n
// returned val does not include \n or \r\n
func read_bytes(bs []byte, size int) (parsed int, val []byte) {
	if len(bs) >= size + 1 {
		if bs[size] == '\r' {
			if len(bs) >= size + 2 {
				if bs[size + 1] == '\n' {
					return size + 2, bs[0 : size]
				} else {
					return -1, nil
				}
			}
		} else if bs[size] == '\n' {
			return size + 1, bs[0 : size]
		} else {
			return -1, nil
		}
	}
	return 0, nil
}
