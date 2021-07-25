// Copyright 2021 The SSDB-cluster Authors
package redis

import (
	"log"
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"base/util"
)

type Request struct {
	Src int
	Dst int
	vals []string
	IsRedis bool
}

func NewRequest(arr []string) *Request {
	ret := new(Request)
	ret.vals = arr
	ret.IsRedis = true
	return ret
}

func DecodeRequest(s string) *Request {
	ret := new(Request)
	ret.Decode([]byte(s))
	return ret
}

func (m *Request)String() string {
	return fmt.Sprintf("%v", m.vals)
}

func (m *Request)Array() []string {
	return m.vals
}

func (m *Request)Cmd() string {
	if len(m.vals) > 0 {
		return m.vals[0]
	}
	return ""
}

func (m *Request)Key() string {
	if len(m.vals) > 1 {
		return m.vals[1]
	}
	return ""
}

func (m *Request)Val() string {
	if len(m.vals) > 2 {
		return m.vals[2]
	}
	return ""
}

func (m *Request)Args() []string {
	if len(m.vals) > 0 {
		return m.vals[1 : ]
	}
	return []string{}
}

func (m *Request)Arg(n int) string {
	if len(m.vals) > 1 + n {
		return m.vals[1 + n]
	}
	return ""
}

func (m *Request)Encode() string {
	if m.IsRedis {
		return m.EncodeRedis()
	} else {
		return m.EncodeSSDB()
	}
}

func (m *Request)EncodeSSDB() string {
	return EncodeSSDB(m.vals)
}

func (m *Request)EncodeRedis() string {
	buf := bytes.NewBuffer(make([]byte, 0, 1 * 1024))
	count := len(m.vals)
	buf.WriteString("*")
	buf.WriteString(strconv.Itoa(count))
	buf.WriteString("\r\n")
	for _, p := range m.vals {
		buf.WriteString("$")
		buf.WriteString(strconv.Itoa(len(p)))
		buf.WriteString("\r\n")
		buf.WriteString(p)
		buf.WriteString("\r\n")
	}
	return buf.String()
}

func (msg *Request)Decode(bs []byte) int {
	// skip leading white spaces
	s := ltrim(bs)
	if s == len(bs) {
		return 0
	}

	var parsed int = 0
	msg.vals = make([]string, 0)

	if bs[s] >= '0' && bs[s] <= '9' {
		msg.vals, parsed = DecodeSSDB(bs[s:])
		msg.IsRedis = false
	} else if bs[s] == '*' || bs[s] == '$' {
		parsed = msg.parseRedisRequest(bs[s:])
		msg.IsRedis = true
	} else {
		end := bytes.IndexByte(bs[s:], '\n')
		if end == -1 {
			return 0
		}
		parsed = end + 1
		if end > s && bs[end-1] == '\r' {
			end -= 1
		}

		msg.vals = util.ParseCommandLine(string(bs[s:s+end]))
		msg.IsRedis = true
	}

	// cmd always lowercased
	if len(msg.vals) > 0 {
		msg.vals[0] = strings.ToLower(msg.vals[0])
	}

	if parsed == -1 {
		return -1
	}
	return s + parsed
}

func (msg *Request)parseRedisRequest(bs []byte) int {
	if len(bs) < 2 {
		return 0
	}

	const BULK  = 0;
	const ARRAY = 1;

	type_ := ARRAY
	bulks := 0

	if (bs[0] == '*') {
		type_  = ARRAY;
		bulks  = 0;
	} else if (bs[0] == '$') {
		type_  = BULK;
		bulks  = 1;
	}

	total := len(bs)

	s := 0
	for s < total {
		if type_ == ARRAY {
			if bs[s] != '*' {
				// log.Println("")
				return -1
			}
		} else if bs[s] != '$' {
			// log.Println("")
			return -1
		}
		s += 1

		idx := bytes.IndexByte(bs[s:], '\n')
		if idx == -1 {
			break
		}
		p := bs[s : s+idx]
		if len(p) > 0 && p[len(p)-1] == '\r' {
			p = p[0 : len(p)-1]
		}
		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			log.Println(err)
			return -1
		}
		s += idx + 1

		if (type_ == ARRAY) {
			bulks  = size
			type_  = BULK
			if bulks == 0 {
				// log.Println("")
				return s
			}
			continue
		}

		end := s + size
		if end >= total { // not ready
			break
		}
		if bs[end] == '\r' {
			end += 1
			if end >= total { // not ready
				break
			}
		}
		if bs[end] != '\n' {
			return -1
		} else {
			p := string(bs[s : s + size])
			msg.vals = append(msg.vals, p)
		}

		s = end + 1
		bulks --
		if bulks == 0 {
			return s
		}
	}

	return 0
}
