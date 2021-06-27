package util

import (
	"strings"
	// "fmt"
)

func is_space(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}

func escape_char(c byte) byte {
	var d byte
	switch c {
	case 'r':
		d = '\r'
	case 'n':
		d = '\n'
	default:
		d = c
	}
	return d
}

func ParseCommandLine(str string) (arr []string) {
	str = strings.TrimSpace(str)
	bs := []byte(str)

	end := len(bs)

	state := 0 // 0: none, 1: single quote, 2: double quote, 3: arg, 4: arg end
	var buf strings.Builder
	var escape bool
	for idx := 0; idx <= end; idx ++ {
		if idx == end && (state == 1 || state == 2) {
			return nil
		}
		if idx == end || state == 4 {
			arr = append(arr, buf.String())
			buf.Reset()
			state = 0
			// fmt.Println("end", idx, end)
		}
		if idx == end {
			break
		}

		c := bs[idx]

		if state == 0 {
			if is_space(c) {
				continue
			} else {
				if c == '\'' {
					state = 1
					continue
				} else if c == '"' {
					state = 2
					continue
				} else {
					state = 3
				}
			}
		}

		if escape {
			c = escape_char(c)
			escape = false
		} else if c == '\\' {
			if state != 0 {
				escape = true
				continue
			}
		} else if is_space(c) {
			if state == 3 {
				state = 4
				continue
			}
		} else if c == '\'' {
			if state == 1 {
				if idx == end - 1 || is_space(bs[idx + 1]) {
					state = 4
					continue
				} else {
					return nil
				}
			} else if state == 3 {
				return nil
			}
		} else if c == '"' {
			if state == 2 {
				if idx == end - 1 || is_space(bs[idx + 1]) {
					state = 4
					continue
				} else {
					return nil
				}
			} else if state == 3 {
				return nil
			}
		}
		// fmt.Printf("put '%c'\n", c)
		buf.WriteByte(c)
	}

	return arr
}