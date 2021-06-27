package util

import (
	"time"
	"math/rand"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

