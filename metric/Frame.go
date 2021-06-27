package metric

type Frame struct {
	stime int64
	etime int64

	count int64
	sum int64
	min int64
	max int64
}

func (g *Frame)Add(latency int64) {
	g.count += 1
	g.sum += latency
	if g.min == 0 || g.min > latency {
		g.min = latency
	}
	if g.max == 0 || g.max < latency {
		g.max = latency
	}
}

func (g *Frame)Duration() int64 {
	return g.etime - g.stime + 1
}
