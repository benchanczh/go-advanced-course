package main

import (
	"sync"
	"time"
)

type numberBucket struct {
	Value float64
}
type Number struct {
	Buckets       map[int64]*numberBucket
	Mutex         *sync.RWMutex
	RollingWindow int64
}

func NewNumber(rollingWindow int64) *Number {
	n := &Number{
		Buckets:       make(map[int64]*numberBucket),
		Mutex:         &sync.RWMutex{},
		RollingWindow: rollingWindow,
	}

	return n
}

func (n *Number) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool
	if bucket, ok = n.Buckets[now]; !ok {
		bucket = &numberBucket{}
		n.Buckets[now] = bucket
	}

	return bucket
}

func (n *Number) removeOldBuckets() {
	now := time.Now().Unix() - n.RollingWindow

	for timestamp := range n.Buckets {
		if timestamp <= now {
			delete(n.Buckets, timestamp)
		}
	}
}

func (n *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	bucket := n.getCurrentBucket()
	bucket.Value += i
	n.removeOldBuckets()
}

func (n *Number) Sum(now time.Time) float64 {
	sum := float64(0)

	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	for timestamp, bucket := range n.Buckets {
		if timestamp >= now.Unix()-n.RollingWindow {
			sum += bucket.Value
		}
	}

	return sum
}

func (n *Number) Max(now time.Time) float64 {
	var max float64

	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	for timestamp, bucket := range n.Buckets {
		if timestamp >= now.Unix()-n.RollingWindow {
			if bucket.Value > max {
				max = bucket.Value
			}
		}
	}

	return max
}

func (n *Number) Average(now time.Time) float64 {
	return n.Sum(now) / float64(n.RollingWindow)
}
