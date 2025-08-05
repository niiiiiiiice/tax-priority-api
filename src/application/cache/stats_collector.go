package cache

import (
	"sync"
	"time"
)

type Stats struct {
	Hits        int64
	Misses      int64
	Sets        int64
	Deletes     int64
	Errors      int64
	LastUpdated time.Time
}

type StatsCollector interface {
	RecordHit()
	RecordMiss()
	RecordError()
	RecordSet()
	RecordDelete()
	GetStats() *Stats
}

type DefaultStatsCollector struct {
	stats   *Stats
	enabled bool
	mu      sync.RWMutex
}

func NewStatsCollector(enabled bool) StatsCollector {
	return &DefaultStatsCollector{
		stats:   &Stats{LastUpdated: time.Now()},
		enabled: enabled,
	}
}

func (s *DefaultStatsCollector) RecordHit() {
	if !s.enabled {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.Hits++
}

func (s *DefaultStatsCollector) RecordMiss() {
	if !s.enabled {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.Misses++
}

func (s *DefaultStatsCollector) RecordError() {
	if !s.enabled {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.Errors++
}

func (s *DefaultStatsCollector) RecordSet() {
	if !s.enabled {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.Sets++
}

func (s *DefaultStatsCollector) RecordDelete() {
	if !s.enabled {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.Deletes++
}

func (s *DefaultStatsCollector) GetStats() *Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statsCopy := *s.stats
	statsCopy.LastUpdated = time.Now()
	return &statsCopy
}
