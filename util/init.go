// Copyright 2021 The SSDB-cluster Authors
package util

import (
	"time"
	"math/rand"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

