// Copyright 2021 The SSDB-cluster Authors
package redis

func ltrim(bs []byte) int {
	s := 0
	for s < len(bs) {
		if bs[s] == ' ' || bs[s] == '\t' || bs[s] == '\r' || bs[s] == '\n' {
			s ++
		} else {
			break
		}
	}
	return s
}
