package cache

import "sync/atomic"

type statsAccessor interface {
	HitCount() uint64
	MissCount() uint64
	LookupCount() uint64
	HitRate() float64
}

type Stats struct {
	hitCount  uint64
	missCount uint64
}

func (s *Stats) IncrHitCount() uint64 {
	return atomic.AddUint64(&s.hitCount, 1)
}

func (s *Stats) IncrMissCount() uint64 {
	return atomic.AddUint64(&s.missCount, 1)
}

func (s *Stats) HitCount() uint64 {
	return atomic.LoadUint64(&s.hitCount)
}

func (s *Stats) MissCount() uint64 {
	return atomic.LoadUint64(&s.missCount)
}

func (s *Stats) LookupCount() uint64 {
	return s.MissCount() + s.HitCount()
}

func (s *Stats) HitRate() float64 {
	hc, mc := s.HitCount(), s.MissCount()
	total := hc + mc
	if total == 0 {
		return 0.0
	}
	return float64(hc) / float64(total)
}
