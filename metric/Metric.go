package metric

import (
	// "fmt"
	"math"
	"sort"
	"time"
	"sync"
)

var FrameDuration int64 = 1000 * 1000

type Metric struct {
	r_mux sync.RWMutex

	frames []*Frame
	max_frames int
}

func NewMetric(max_frames int) *Metric {
	m := new(Metric)
	m.max_frames = max_frames
	return m
}

func (m *Metric)Add(latency int64) {
	now := time.Now().UnixNano() / 1000

	m.r_mux.Lock()
	defer m.r_mux.Unlock()

	if len(m.frames) == 0 {
		g := new(Frame)
		g.stime = now
		m.frames = append(m.frames, g)
	}
	if len(m.frames) >= m.max_frames * 2 {
		copy(m.frames, m.frames[m.max_frames : ])
		m.frames = m.frames[0 : m.max_frames]
	}

	curr := m.frames[len(m.frames) - 1]

	if now - curr.stime + 1 > FrameDuration {
		curr.etime = now - 1
		// fmt.Printf("%d %d\n", curr.stime, curr.etime)

		curr = new(Frame)
		curr.stime = now
		m.frames = append(m.frames, curr)
	}

	curr.etime = now
	curr.Add(latency)
}

func (m *Metric)CalcUntilNow() *Stat {
	res := new(Stat)
	var frames []*Frame

	m.r_mux.RLock()
	if len(m.frames) > 0 {
		frames = make([]*Frame, len(m.frames))
		copy(frames, m.frames)
	}
	m.r_mux.RUnlock()

	if len(frames) > 0 {
		dummy := new(Frame)
		dummy.stime = frames[len(frames) - 1].etime + 1
		dummy.etime = time.Now().UnixNano() / 1000

		frames = append(frames, dummy)
		calc(frames, res)
	}

	return res
}

func (m *Metric)Calc(stime, etime int64) *Stat {
	res := new(Stat)
	var frames []*Frame

	m.r_mux.RLock()
	spos := sort.Search(len(m.frames), func(i int) bool {
		return m.frames[i].etime >= stime
	})
	epos := sort.Search(len(m.frames), func(i int) bool {
		return m.frames[i].stime > etime
	})
	if epos - spos > 0 {
		frames = make([]*Frame, epos - spos)
		copy(frames, m.frames[spos : epos])

		if epos == len(m.frames) {
			// last sample is mutable, clone it
			var ns Frame
			ns = *(frames[len(frames)-1])
			frames[len(frames)-1] = &ns
		}
	}
	m.r_mux.RUnlock()

	if len(frames) > 0 {
		for i, s := range frames {
			if i > 0 && i < len(frames) - 1 {
				continue
			}

			count := s.count
			sum := s.sum
			duration := float64(s.Duration())
			{
				r := float64(s.etime - stime + 1) / duration
				r = limit(r, 0, 1)
				count = int64(math.Round(float64(count) * r))
				sum = int64(math.Round(float64(sum) * r))
			}
			{
				r := float64(etime - s.stime + 1) / duration
				r = limit(r, 0, 1)
				count = int64(math.Round(float64(count) * r))
				sum = int64(math.Round(float64(sum) * r))
			}

			n := new(Frame)
			if i == 0 {
				n.stime = stime
			} else {
				n.stime = s.stime
			}
			if i == len(frames) - 1 {
				n.etime = etime
			} else {
				n.etime = s.etime
			}
			n.count = count
			n.sum = sum
			n.min = s.min
			n.max = s.max
			frames[i] = n
		}

		// fmt.Println("source:")
		// for i, s := range m.frames {
		// 	fmt.Printf("  [%2d] %d %d\n", i, s.stime, s.etime)
		// }
		// fmt.Printf("input: %d %d\n", stime, etime)
		// fmt.Println("output:")
		// for i, s := range frames {
		// 	fmt.Printf("  [%2d] %d %d\n", i, s.stime, s.etime)
		// }

		calc(frames, res)
	}

	return res
}

func calc(frames []*Frame, res *Stat) {
	var latency_sum int64
	for _, g := range frames {
		res.Count += g.count
		if res.Min == 0 || (g.min > 0 && g.min < res.Min) {
			res.Min = g.min
		}
		if res.Max == 0 || (g.max > 0 && g.max > res.Max) {
			res.Max = g.max
		}
		latency_sum += g.sum
	}

	res.Time = frames[len(frames) - 1].etime - frames[0].stime + 1
	if res.Time > 0 {
		res.Qps = 1000 * 1000 * res.Count / res.Time
	}
	if res.Count > 0 {
		res.Avg = latency_sum / res.Count
	}
}

func limit(n, min, max float64) float64 {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
