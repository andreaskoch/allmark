// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Code for ViewModelCache created
// from https://github.com/streamrail/concurrent-map/blob/master/concurrent_map_template.txt

package orchestrator

import (
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"hash/fnv"
	"sync"
)

var VIEWMODELCACHE_SHARD_COUNT = 32

// A "thread" safe map of type string:viewmodel.Model.
// To avoid lock bottlenecks this map is dived to several (VIEWMODELCACHE_SHARD_COUNT) map shards.
type ViewModelCache []*ConcurrentViewModelMapShared
type ConcurrentViewModelMapShared struct {
	items        map[string]viewmodel.Model
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent viewmodel cache map.
func newViewmodelCache() ViewModelCache {
	m := make(ViewModelCache, VIEWMODELCACHE_SHARD_COUNT)
	for i := 0; i < VIEWMODELCACHE_SHARD_COUNT; i++ {
		m[i] = &ConcurrentViewModelMapShared{items: make(map[string]viewmodel.Model)}
	}
	return m
}

// Returns shard under given key
func (m ViewModelCache) GetShard(key string) *ConcurrentViewModelMapShared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[hasher.Sum32()%uint32(VIEWMODELCACHE_SHARD_COUNT)]
}

// Sets the given value under the specified key.
func (m *ViewModelCache) Set(key string, value viewmodel.Model) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// Retrieves an element from map under given key.
func (m ViewModelCache) Get(key string) (viewmodel.Model, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Returns the number of elements within the map.
func (m ViewModelCache) Count() int {
	count := 0
	for i := 0; i < VIEWMODELCACHE_SHARD_COUNT; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Looks up an item under specified key
func (m *ViewModelCache) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Removes an element from the map.
func (m *ViewModelCache) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// Checks if map is empty.
func (m *ViewModelCache) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type ViewModelCacheTuple struct {
	Key string
	Val viewmodel.Model
}

// Returns an iterator which could be used in a for range loop.
func (m ViewModelCache) Iter() <-chan ViewModelCacheTuple {
	ch := make(chan ViewModelCacheTuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ViewModelCacheTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// Returns a buffered iterator which could be used in a for range loop.
func (m ViewModelCache) IterBuffered() <-chan ViewModelCacheTuple {
	ch := make(chan ViewModelCacheTuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ViewModelCacheTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}
