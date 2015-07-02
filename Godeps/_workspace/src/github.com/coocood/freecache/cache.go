package freecache

import (
	"crypto/md5"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Cache struct {
	locks     [256]sync.Mutex
	segments  [256]segment
	hitCount  int64
	missCount int64
}

func hashFunc(data []byte) uint64 {
	sum := md5.Sum(data)
	return *(*uint64)(unsafe.Pointer(&sum[0]))
}

// The cache size will be set to 512KB at minimum.
// If the size is set relatively large, you should call
// `debug.SetGCPercent()`, set it to a much smaller value
// to limit the memory consumption and GC pause time.
func NewCache(size int) (cache *Cache) {
	if size < 512*1024 {
		size = 512 * 1024
	}
	cache = new(Cache)
	for i := 0; i < 256; i++ {
		cache.segments[i] = newSegment(size/256, i)
	}
	return
}

// If the key is larger than 65535 or value is larger than 1/1024 of the cache size,
// the entry will not be written to the cache. expireSeconds <= 0 means no expire,
// but it can be evicted when cache is full.
func (cache *Cache) Set(key, value []byte, expireSeconds int) (err error) {
	hashVal := hashFunc(key)
	segId := hashVal & 255
	cache.locks[segId].Lock()
	err = cache.segments[segId].set(key, value, hashVal, expireSeconds)
	cache.locks[segId].Unlock()
	return
}

// Get the value or not found error.
func (cache *Cache) Get(key []byte) (value []byte, err error) {
	hashVal := hashFunc(key)
	segId := hashVal & 255
	cache.locks[segId].Lock()
	value, err = cache.segments[segId].get(key, hashVal)
	cache.locks[segId].Unlock()
	if err == nil {
		atomic.AddInt64(&cache.hitCount, 1)
	} else {
		atomic.AddInt64(&cache.missCount, 1)
	}
	return
}

func (cache *Cache) Del(key []byte) (affected bool) {
	hashVal := hashFunc(key)
	segId := hashVal & 255
	cache.locks[segId].Lock()
	affected = cache.segments[segId].del(key, hashVal)
	cache.locks[segId].Unlock()
	return
}

func (cache *Cache) EvacuateCount() (count int64) {
	for i := 0; i < 256; i++ {
		count += atomic.LoadInt64(&cache.segments[i].totalEvacuate)
	}
	return
}

func (cache *Cache) EntryCount() (entryCount int64) {
	for i := 0; i < 256; i++ {
		entryCount += atomic.LoadInt64(&cache.segments[i].entryCount)
	}
	return
}

// The average unix timestamp when a entry being accessed.
// Entries have greater access time will be evacuated when it
// is about to be overwritten by new value.
func (cache *Cache) AverageAccessTime() int64 {
	var entryCount, totalTime int64
	for i := 0; i < 256; i++ {
		totalTime += atomic.LoadInt64(&cache.segments[i].totalTime)
		entryCount += atomic.LoadInt64(&cache.segments[i].totalCount)
	}
	if entryCount == 0 {
		return 0
	} else {
		return totalTime / entryCount
	}
}

func (cache *Cache) HitCount() int64 {
	return atomic.LoadInt64(&cache.hitCount)
}

func (cache *Cache) LookupCount() int64 {
	return atomic.LoadInt64(&cache.hitCount) + atomic.LoadInt64(&cache.missCount)
}

func (cache *Cache) HitRate() float64 {
	lookupCount := cache.LookupCount()
	if lookupCount == 0 {
		return 0
	} else {
		return float64(cache.HitCount()) / float64(lookupCount)
	}
}

func (cache *Cache) OverwriteCount() (overwriteCount int64) {
	for i := 0; i < 256; i++ {
		overwriteCount += atomic.LoadInt64(&cache.segments[i].overwrites)
	}
	return
}

func (cache *Cache) Clear() {
	for i := 0; i < 256; i++ {
		cache.locks[i].Lock()
		newSeg := newSegment(len(cache.segments[i].rb.data), i)
		cache.segments[i] = newSeg
		cache.locks[i].Unlock()
	}
	atomic.StoreInt64(&cache.hitCount, 0)
	atomic.StoreInt64(&cache.missCount, 0)
}