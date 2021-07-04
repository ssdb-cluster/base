package redis

import (
	"bytes"
	"strconv"
)

const (
	TypeOK = iota
	TypeError
	TypeNull
	TypeInt
	TypeString
	TypeArray
)

// simple reply, do not support nested replies
type Response struct {
	Dst int
	type_ int
	vals []string
}

func (r *Response)ErrorCode() string {
	if r.type_ == TypeError {
		if len(r.vals) > 0 {
			return r.vals[0]
		}
	}
	return ""
}

func (r *Response)ErrorMessage() string {
	if r.type_ == TypeError {
		if len(r.vals) > 1 {
			return r.vals[1]
		}
	}
	return ""
}

// if reply type is not INT, return -1
func (r *Response)Int() int {
	if r.IsInt() && len(r.vals) > 0 {
		i, _ := strconv.Atoi(r.vals[0])
		return i
	}
	return -1
}

// if reply type is not STRING, return ""
func (r *Response)String() string {
	if r.IsString() && len(r.vals) > 0 {
		return r.vals[0]
	}
	return ""
}

func (r *Response)Array() []string {
	return r.vals
}

func (r *Response)Kvs() (kvs [][2]string) {
	for i := 0; i < len(r.vals) - 1; i += 2 {
		kvs = append(kvs, [2]string{r.vals[i], r.vals[i+1]})
	}
	return
}

func (r *Response)Pairs() (kvs [][2]string) {
	return r.Kvs()
}

func (r *Response)IsOK() bool {
	return r.type_ == TypeOK
}

func (r *Response)IsError() bool {
	return r.type_ == TypeError
}

func (r *Response)IsNull() bool {
	return r.type_ == TypeNull
}

func (r *Response)IsInt() bool {
	return r.type_ == TypeInt
}

func (r *Response)IsString() bool {
	return r.type_ == TypeString
}

func (r *Response)IsArray() bool {
	return r.type_ == TypeArray
}

func (r *Response)SetError(msg string) {
	r.type_ = TypeError
	r.vals = []string{"ERR", msg}
}

func (r *Response)SetError2(code string, msg string) {
	r.type_ = TypeError
	r.vals = []string{code, msg}
}

func (r *Response)SetNull() {
	r.type_ = TypeNull
}

func (r *Response)SetInt(num int64) {
	r.type_ = TypeInt
	r.vals = []string{strconv.FormatInt(num, 10)}
}

func (r *Response)SetString(b string) {
	r.type_ = TypeString
	r.vals = []string{b}
}

func (r *Response)SetArray(ps []string) {
	r.type_ = TypeArray
	r.vals = ps
}

func (r *Response)EncodeSSDB() string {
	vals := r.vals
	switch r.type_ {
	case TypeOK:
		vals = []string{"ok"}
	case TypeNull:
		vals = []string{"not_found"}
	case TypeError:
		vals[0] = "error"
	default:
		vals = append([]string{"ok"}, vals...)
	}
	return EncodeSSDB(vals)
}

func (r *Response)Encode() string {
	buf := bytes.NewBuffer(make([]byte, 0, 1 * 1024))
	switch r.type_ {
	case TypeOK:
		return "+OK\r\n";
	case TypeError:
		buf.WriteByte('-');
		buf.WriteString(r.vals[0]);
		buf.WriteString(" ");
		buf.WriteString(r.vals[1]);
		buf.WriteString("\r\n");
	case TypeNull:
		return "$-1\r\n";
	case TypeInt:
		buf.WriteByte(':');
		buf.WriteString(r.vals[0]);
		buf.WriteString("\r\n");
	case TypeString:
		buf.WriteByte('$');
		buf.WriteString(strconv.Itoa(len(r.vals[0])));
		buf.WriteString("\r\n");
		buf.WriteString(r.vals[0]);
		buf.WriteString("\r\n");
	case TypeArray:
		buf.WriteByte('*');
		buf.WriteString(strconv.Itoa(len(r.vals)));
		buf.WriteString("\r\n");
		for _, s := range r.vals {
			buf.WriteByte('$');
			buf.WriteString(strconv.Itoa(len(s)));
			buf.WriteString("\r\n");
			buf.WriteString(s);
			buf.WriteString("\r\n");
		}
	}
	return buf.String()
}

func (r *Response)Decode(bs []byte) int {
	// skip leading white spaces
	s := ltrim(bs)
	if s == len(bs) {
		return 0
	}

	reply := new(Reply)
	parsed := reply.Decode(bs[s :])
	if parsed <= 0 {
		return parsed
	}
	s += parsed

	r.type_ = reply.type_
	r.vals = make([]string, 0)

	if reply.type_ != TypeArray {
		if reply.type_ == TypeError {
			r.vals = append(r.vals, reply.err)
			if len(reply.val) > 0 {
				r.vals = append(r.vals, reply.val)
			}
		} else if reply.type_ == TypeNull {
			//
		} else {
			r.vals = append(r.vals, reply.val)
		}
		return s
	}

	for i := 0; i < reply.num; i ++ {
		reply := new(Reply)
		parsed := reply.Decode(bs[s :])
		if parsed <= 0 {
			return parsed
		}
		s += parsed
		// only support array of string
		if reply.type_ != TypeString {
			return -1
		}
		r.vals = append(r.vals, reply.val)
	}

	return s
}
