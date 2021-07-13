// Copyright 2021 The SSDB-cluster Authors
package jconf

func Decode(s string) (*Object, error) {
	obj := NewObject()
	err := obj.Decode(s)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
