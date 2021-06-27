package metric

import (
	"fmt"
)

type Stat struct {
	Time int64
	Count int64
	Qps int64
	Avg int64
	Min int64
	Max int64
}

func (s *Stat)String() string {
	return fmt.Sprintf(`Stats:
Time:  %.03f s
Count: %d
Qps:   %d
Avg:   %d us
Max:   %d us
Min:   %d us`,
	float64(s.Time)/1000/1000, s.Count, s.Qps, s.Avg, s.Max, s.Min)
}
