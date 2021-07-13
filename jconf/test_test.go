// Copyright 2021 The SSDB-cluster Authors
package jconf

import (
	"testing"
	"log"
	"io/ioutil"
)

func Test_basic(t *testing.T){
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	bs, _ := ioutil.ReadFile("a.json")
	s := string(bs)

	obj, err := Decode(s)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%v", obj.Encode())
}
