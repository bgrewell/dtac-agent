package common

import (
	"math"
	"sync"
	"time"
)

type TimestampedFloatArrayEntry struct {
	Timestamp int64 `json:"timestamp"`
	Value float64 `json:"value"`
}

func (te *TimestampedFloatArrayEntry) Before(timestamp int64) bool {
	return te.Timestamp < timestamp
}

func (te *TimestampedFloatArrayEntry) After(timestamp int64) bool {
	return te.Timestamp > timestamp
}

type TimestampedFloatArray struct {
	Count int `json:"count"`
	Timeouts int `json:"timeouts"`
	Entries []TimestampedFloatArrayEntry `json:"entries"`
	mutex sync.Mutex
}

func (t *TimestampedFloatArray) Add(entry float64) {
	if t.Entries == nil {
		t.mutex.Lock()
		t.Entries = make([]TimestampedFloatArrayEntry, 0)
		t.mutex.Unlock()
	}
	if entry >= 0 {
		e := TimestampedFloatArrayEntry{
			Timestamp: time.Now().UnixNano(),
			Value:     entry,
		}
		t.mutex.Lock()
		t.Count += 1
		t.Entries = append(t.Entries, e)
		t.mutex.Unlock()
	} else {
		// todo: need to figure out a good way to handle timeouts, don't want to skew the results but don't want to ignore. set to timeout value?
		t.mutex.Lock()
		t.Timeouts += 1
		t.mutex.Unlock()
	}
}

func (t *TimestampedFloatArray) FirstNSec(seconds int) []TimestampedFloatArrayEntry {
	entries := make([]TimestampedFloatArrayEntry, 0)
	cutoff := t.Entries[0].Timestamp / int64(time.Second) + int64(seconds)
	for i := 0; i < len(t.Entries); i++ {
		if t.Entries[i].After(cutoff) {
			return entries
		}
		entries = append(entries, t.Entries[i])
	}
	return entries
}

func (t *TimestampedFloatArray) LastNSec(seconds int) []TimestampedFloatArrayEntry {
	entries := make([]TimestampedFloatArrayEntry, 0)
	last := len(t.Entries) - 1
	cutoff := t.Entries[last].Timestamp / int64(time.Second) - int64(seconds)
	for i := last; i >= 0; i-- {
		if t.Entries[i].Before(cutoff) {
			return entries
		}
		entries = append(entries, t.Entries[i])
	}
	return entries
}

func (t *TimestampedFloatArray) Average() float64 {
	total := float64(0)
	num := float64(len(t.Entries))
	for _, entry := range t.Entries {
		total += entry.Value
	}
	return total / num
}

func (t *TimestampedFloatArray) AveragePeriod(seconds int) float64 {
	entries := t.LastNSec(seconds)
	total := float64(0)
	num := float64(len(entries))
	for _, entry := range entries {
		total += entry.Value
	}
	return total / num
}

func (t *TimestampedFloatArray) StdDev() float64 {
	mean := t.Average()
	sd := float64(0)
	num := float64(len(t.Entries))
	for _, entry := range t.Entries {
		sd += math.Pow(entry.Value - mean, 2)
	}
	return math.Sqrt(sd / num)
}

func (t *TimestampedFloatArray) StdDevPeriod(seconds int) float64 {
	entries := t.LastNSec(seconds)
	mean := t.Average()
	sd := float64(0)
	num := float64(len(entries))
	for _, entry := range entries {
		sd += math.Pow(entry.Value - mean, 2)
	}
	return math.Sqrt(sd / num)
}
