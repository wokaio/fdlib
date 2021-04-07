// Copyright (c) 2021 Miczone Asia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package caching

import (
	"container/list"
	"fmt"
	"sync"
)

type LRUCaching struct {
	lock     sync.Mutex
	capacity int                           // maximum number of key-value pairs
	cache    map[interface{}]*list.Element // map for cached key-value pairs
	lru      *list.List                    // LRU list
}

type Pair struct {
	key   interface{} // cache key
	value interface{} // cache value
}

// NewLRUCaching returns a new, empty LRUCaching
func NewLRUCaching(capacity int) *LRUCaching {
	c := new(LRUCaching)
	c.capacity = capacity
	c.cache = make(map[interface{}]*list.Element)
	c.lru = list.New()
	return c
}

// Get get cached value from LRU cache
// The second return value indicates whether key is found or not, true if found, false if not
func (c *LRUCaching) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.lru.MoveToFront(elem) // move node to head of lru list
		return elem.Value.(*Pair).value, true
	}
	return nil, false
}

// Add adds a key-value pair to LRU cache, true if eviction occurs, false if not
func (c *LRUCaching) Add(key interface{}, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.lru.MoveToFront(elem) // update lru list
		elem.Value.(*Pair).value = value
		return false
	}

	elem := c.lru.PushFront(&Pair{key, value})
	c.cache[key] = elem

	if c.lru.Len() > c.capacity {
		c.Evict()
		return true
	}

	return false
}

// Evict a key-value pair from LRU cache
func (c *LRUCaching) Evict() {
	elem := c.lru.Back()
	if elem == nil {
		return
	}
	// remove item at the end of lru list
	c.lru.Remove(elem)
	delete(c.cache, elem.Value.(*Pair).key)
}

// Del deletes cached value from cache
func (c *LRUCaching) Del(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.lru.Remove(elem)
		delete(c.cache, key)
	}
}

// Len returns number of items in cache
func (c *LRUCaching) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Len()
}

// Keys returns keys of items in cache
func (c *LRUCaching) Keys() []interface{} {
	var keyList []interface{}
	c.lock.Lock()
	for key := range c.cache {
		keyList = append(keyList, key)
	}
	c.lock.Unlock()
	return keyList
}

//EnlargeCapacity enlarges the capacity of cache
func (c *LRUCaching) EnlargeCapacity(newCapacity int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if newCapacity < c.capacity {
		return fmt.Errorf("newCapacity[%d] must be larger than currentCapacity[%d]",
			newCapacity, c.capacity)
	}
	c.capacity = newCapacity

	return nil
}
