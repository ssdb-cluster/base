package util

import (
	"fmt"
	"strings"
	"bytes"
)

func ReplaceBytes(s string, src []string, dst []string) string {
	for i, _ := range src {
		s = strings.Replace(s, src[i], dst[i], -1)
	}
	return s
}

// 仅用于调试
func StringEscape(s string) string {
	s = fmt.Sprintf("%q", s)
	return s[1:len(s)-1]
	// return string(BytesEscape([]byte(s)))
}

// 仅用于调试
func StringUnescape(s string) string {
	var ret string
	fmt.Sscanf("\"" + s + "\"", "%q", &ret)
	return ret
	// return string(BytesUnescape([]byte(s)))
}

// See strconv.Quote() https://golang.org/src/strconv/quote.go?s=4976:5005
func BytesEscape(bs []byte) []byte {
	var buf bytes.Buffer
	var s int = 0
	var e int = -1
	var c byte
	for e, c = range bs {
		var d string
		switch c {
		case '\\':
			d = "\\\\"
		case ' ':
			d = "\\s"
		case '\t':
			d = "\\t"
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

func BytesUnescape(bs []byte) []byte {
	var buf bytes.Buffer
	var s int = 0
	var e int = -1
	var p byte = 0
	var c byte = 0
	for e, c = range bs {
		// log.Printf("%c", c)
		if p == '\\' {
			var d byte
			switch c {
			case '\\':
				d = '\\'
			case 's':
				d = ' '
			case 't':
				d = '\t'
			case 'r':
				d = '\r'
			case 'n':
				d = '\n'
			default:
				p = c
				continue
			}
			// log.Println(s, e, len(bs))
			buf.Write(bs[s : e - 1])
			buf.WriteByte(d)
			s = e + 1
		}
		p = c
	}
	if s == 0 && e == len(bs) - 1 {
		return bs // no copy
	}
	if s <= e {
		buf.Write(bs[s : e + 1])
	}
	return buf.Bytes()
}
