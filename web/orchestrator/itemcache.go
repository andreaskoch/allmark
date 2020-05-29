// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Code for ItemCache created
// from https://github.com/streamrail/concurrent-map/blob/master/concurrent_map_template.txt

package orchestrator

import (
	"github.com/andreaskoch/allmark/model"
	"hash/fnv"
	"sync"
)

var ITEMCACHE_SHARD_COUNT = uint(32)

// A "thread" safe map of type string:*model.Item.
// To avoid lock bottlenecks this map is dived to several (ITEMCACHE_SHARD_COUNT) map shards.
type ItemCache []*ConcurrentItemMapShared
type ConcurrentItemMapShared struct {
	items        map[string]*model.Item
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent item cache map.
func newItemCache() ItemCache {
	m := make(ItemCache, ITEMCACHE_SHARD_COUNT)
	for i := uint(0); i < ITEMCACHE_SHARD_COUNT; i++ {
		m[i] = &ConcurrentItemMapShared{items: make(map[string]*model.Item)}
	}
	return m
}

// Returns shard under given key
func (m ItemCache) GetShard(key string) *ConcurrentItemMapShared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[uint(hasher.Sum32())%ITEMCACHE_SHARD_COUNT]
}

// Sets the given value under the specified key.
func (m *ItemCache) Set(key string, value *model.Item) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// Retrieves an element from map under given key.
func (m ItemCache) Get(key string) (*model.Item, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Returns the number of elements within the map.
func (m ItemCache) Count() int {
	count := 0
	for i := uint(0); i < ITEMCACHE_SHARD_COUNT; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Looks up an item under specified key
func (m *ItemCache) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Removes an element from the map.
func (m *ItemCache) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// Checks if map is empty.
func (m *ItemCache) IsEmpty() bool {
	return m.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type ItemCacheTuple struct {
	Key string
	Val *model.Item
}

// Returns an iterator which could be used in a for range loop.
func (m ItemCache) Iter() <-chan ItemCacheTuple {
	ch := make(chan ItemCacheTuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ItemCacheTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// Returns a buffered iterator which could be used in a for range loop.
func (m ItemCache) IterBuffered() <-chan ItemCacheTuple {
	ch := make(chan ItemCacheTuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ItemCacheTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}
