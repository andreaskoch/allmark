// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Code for ViewModelListCache created
// from https://github.com/streamrail/concurrent-map/blob/master/concurrent_map_template.txt

package orchestrator

import (
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"hash/fnv"
	"sync"
)

var VIEWMODELLISTCACHE_SHARD_COUNT = uint(32)

// A "thread" safe map of type string:[]viewmodel.Model.
// To avoid lock bottlenecks this map is dived to several (VIEWMODELLISTCACHE_SHARD_COUNT) map shards.
type ViewModelListCache []*ConcurrentViewModelListMapShared
type ConcurrentViewModelListMapShared struct {
	items        map[string][]viewmodel.Model
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent viewmodel cache map.
func newViewModelListCache() ViewModelListCache {
	m := make(ViewModelListCache, VIEWMODELLISTCACHE_SHARD_COUNT)
	for i := uint(0); i < VIEWMODELLISTCACHE_SHARD_COUNT; i++ {
		m[i] = &ConcurrentViewModelListMapShared{items: make(map[string][]viewmodel.Model)}
	}
	return m
}

// Returns shard under given key
func (m ViewModelListCache) GetShard(key string) *ConcurrentViewModelListMapShared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[uint(hasher.Sum32())%VIEWMODELLISTCACHE_SHARD_COUNT]
}

// Sets the given value under the specified key.
func (m *ViewModelListCache) Set(key string, value []viewmodel.Model) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// Retrieves an element from map under given key.
func (m ViewModelListCache) Get(key string) ([]viewmodel.Model, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Returns the number of elements within the map.
func (m ViewModelListCache) Count() int {
	count := 0
	for i := uint(0); i < VIEWMODELLISTCACHE_SHARD_COUNT; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Looks up an item under specified key
func (m *ViewModelListCache) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Removes an element from the map.
func (m *ViewModelListCache) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// Checks if map is empty.
func (m *ViewModelListCache) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type ViewModelListCacheTuple struct {
	Key string
	Val []viewmodel.Model
}

// Returns an iterator which could be used in a for range loop.
func (m ViewModelListCache) Iter() <-chan ViewModelListCacheTuple {
	ch := make(chan ViewModelListCacheTuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ViewModelListCacheTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// Returns a buffered iterator which could be used in a for range loop.
func (m ViewModelListCache) IterBuffered() <-chan ViewModelListCacheTuple {
	ch := make(chan ViewModelListCacheTuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ViewModelListCacheTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}
