package redis

import (
	"bytes"
	"strconv"
)

// TODO: return []byte
func EncodeSSDB(arr []string) string {
	if len(arr) == 0 {
		return "0\n\n\n"
	}
	buf := bytes.NewBuffer(make([]byte, 0, 1 * 1024))
	for _, p := range arr {
		buf.WriteString(strconv.Itoa(len(p)))
		buf.WriteByte('\n')
		buf.WriteString(p)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	return buf.String()
}

func DecodeSSDB(bs []byte) (arr []string, nn int) {
	s := 0
	total := len(bs)

	for {
		idx := bytes.IndexByte(bs[s:], '\n')
		if idx == -1 {
			break
		}

		p := bs[s : s+idx]
		s += idx + 1
		if len(p) > 0 && p[0] == '\r' {
			p = p[0 : len(p)-1]
		}
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			// log.Printf("parse end")
			return arr, s
		}
		// log.Printf("> size [%s]\n", p);

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil, -1
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
			return nil, -1
		} else {
			p := string(bs[s : s + size])
			arr = append(arr, p)
			s = end + 1
			// log.Printf("> data %d %d [%s]\n", start, size, p);
		}
	}
	return nil, 0
}
