package jconf

import (
	"fmt"
	"strconv"
	"bytes"
	"strings"
)

const (
	OBJ_NULL   = 0
	OBJ_STRING = 1
	OBJ_ARRAY  = 2
	OBJ_OBJECT = 3
)

var NullObj = &Object{type_: OBJ_NULL, value: "null"}

type Object struct {
	type_ int
	key string
	value string
	subs []*Object
}

func NewString(s string) *Object {
	ret := new(Object)
	ret.type_ = OBJ_STRING
	ret.value = s
	return ret
}

func NewArray() *Object {
	ret := new(Object)
	ret.type_ = OBJ_ARRAY
	return ret
}

func NewObject() *Object {
	ret := new(Object)
	ret.type_ = OBJ_OBJECT
	return ret
}

func (obj *Object)IsNull() bool {
	return obj.type_ == OBJ_NULL
}

func (obj *Object)HasKey(key string) bool {
	return !obj.Get(key).IsNull()
}

func (obj *Object)Count() int {
	return len(obj.subs)
}

func (obj *Object)Empty() bool {
	if obj.type_ == OBJ_STRING {
		return len(obj.value) == 0
	}
	return len(obj.subs) == 0
}

func (obj *Object)String() (ret string) {
	switch obj.type_ {
	case OBJ_OBJECT:
		ret = "[object]"
	case OBJ_ARRAY:
		ret = "[array]"
	default:
		ret = obj.value
	}
	return
}

func (obj *Object)Integer() int64 {
	r, _ := strconv.ParseInt(obj.value, 10, 64)
	return r
}

func (obj *Object)Array() []*Object {
	return obj.subs
}

func (obj *Object)Keys() (arr []string) {
	for i, sub := range obj.subs {
		if obj.type_ == OBJ_ARRAY {
			arr = append(arr, fmt.Sprintf("%d", i))
		} else {
			arr = append(arr, sub.key)
		}
	}
	return
}

// return all subs as string array
func (obj *Object)Values() (arr []string) {
	for _, sub := range obj.subs {
		arr = append(arr, sub.String())
	}
	return
}

func interface_to_obj(val interface{}) *Object {
	var sub *Object
	if _, ok := val.(*Object); ok {
		sub = val.(*Object)
	} else {
		sub = NewString(fmt.Sprintf("%v", val))
	}
	return sub
}

// object method
func (obj *Object)Set(key string, val interface{}) {
	if obj.IsNull() {
		return
	}

	sub := interface_to_obj(val)
	sub.key = key

	obj.type_ = OBJ_OBJECT

	for i, s := range obj.subs {
		if s.key == key {
			obj.subs[i] = sub
			return
		}
	}

	obj.subs = append(obj.subs, sub)
}

// object method
func (obj *Object)Unset(key string) {
	for i, s := range obj.subs {
		if s.key == key {
			obj.subs = append(obj.subs[0:i], obj.subs[i+1:]...)
			return
		}
	}
}

// object method
func (obj *Object)Get(key string) *Object {
	for _, s := range obj.subs {
		if s.key == key {
			return s
		}
	}
	return NullObj
}

// array method
func (obj *Object)ItemAt(idx int) *Object {
	if idx > len(obj.subs) {
		return NullObj
	}
	return obj.subs[idx]
}

// object|array method
func (obj *Object)Clear() {
	if obj.IsNull() {
		return
	}
	obj.subs = make([]*Object, 0)
}

// array method
func (obj *Object)Push(val interface{}) {
	if obj.IsNull() {
		return
	}

	sub := interface_to_obj(val)

	obj.type_ = OBJ_ARRAY
	obj.subs = append(obj.subs, sub)
}

// array method
func (obj *Object)Pop() *Object {
	if obj.IsNull() {
		return nil
	}

	ret := obj.subs[len(obj.subs) - 1]
	obj.type_ = OBJ_ARRAY
	obj.subs = obj.subs[0 : len(obj.subs) - 1]
	return ret
}

func (obj *Object)Decode(s string) error {
	p := &Parser{tn: &Tokenizer{buf: []byte(s)}}
	ret, err := p.parse()
	if err != nil {
		return err
	}
	obj.type_ = ret.type_
	obj.value = ret.value
	obj.subs = ret.subs
	return nil
}

func (obj *Object)Encode() string {
	buf := new(bytes.Buffer)
	obj.encode_to_buffer(buf, 0)
	buf.WriteByte('\n')
	return buf.String()
}

func (obj *Object)is_string_array() bool {
	for _, sub := range obj.subs {
		if sub.type_ != OBJ_STRING {
			return false
		}
	}
	return true
}

func (obj *Object)encode_to_buffer(buf *bytes.Buffer, depth int) {
	switch obj.type_ {
	case OBJ_OBJECT:
		if len(obj.subs) == 0 {
			buf.WriteString("{}")
		} else {
			indent := strings.Repeat("    ", depth)
			sub_indent := strings.Repeat("    ", depth+1)
			buf.WriteString("{\n")
			for i, sub := range obj.subs {
				if i != 0 {
					buf.WriteString(",\n")
				}
				buf.WriteString(sub_indent)
				buf.WriteString("\"")
				buf.WriteString(sub.key)
				buf.WriteString("\": ")
				sub.encode_to_buffer(buf, depth + 1)
			}
			buf.WriteString("\n")
			buf.WriteString(indent)
			buf.WriteString("}")
		}
	case OBJ_ARRAY:
		if len(obj.subs) == 0 {
			buf.WriteString("[]")
		} else {
			if obj.is_string_array() {
				length := 2 // []
				for _, sub := range obj.subs {
					length += len(sub.value) + 4 // "",\s
				}
				if length < 80 {
					buf.WriteString("[")
					for i, sub := range obj.subs {
						if i != 0 {
							buf.WriteString(", ")
						}
						sub.encode_to_buffer(buf, depth + 1)
					}
					buf.WriteString("]")
					return
				}
			}

			indent := strings.Repeat("    ", depth)
			sub_indent := strings.Repeat("    ", depth+1)
			buf.WriteString("[\n")
			for i, sub := range obj.subs {
				if i != 0 {
					buf.WriteString(",\n")
				}
				buf.WriteString(sub_indent)
				sub.encode_to_buffer(buf, depth + 1)
			}
			buf.WriteString("\n")
			buf.WriteString(indent)
			buf.WriteString("]")
		}
	default:
		fmt.Fprintf(buf, "%q", obj.value)
	}
}

func (obj *Object)Clone() *Object {
	ret := new(Object)
	ret.type_ = obj.type_
	ret.key = obj.key
	ret.value = obj.value
	ret.subs = make([]*Object, len(obj.subs))
	for i, sub := range obj.subs {
		ret.subs[i] = sub.Clone()
	}
	return ret
}
