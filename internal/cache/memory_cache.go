/*
 * Copyright 2025 Praveen Kumar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import (
	"sync"
	"time"
)

type InMemoryCache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]*CacheEntry[V]
}

func NewInMemoryCache[K comparable, V any]() *InMemoryCache[K, V] {
	return &InMemoryCache[K, V]{
		items: make(map[K]*CacheEntry[V]),
	}
}

func (c *InMemoryCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[key]
	if !exists {
		var zero V
		return zero, false
	}

	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		var zero V
		return zero, false
	}

	return entry.Data, true
}

func (c *InMemoryCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheEntry[V]{
		Data:     value,
		CachedAt: time.Now(),
	}
}

func (c *InMemoryCache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	c.items[key] = &CacheEntry[V]{
		Data:      value,
		CachedAt:  time.Now(),
		ExpiresAt: &expiresAt,
	}
}

func (c *InMemoryCache[K, V]) Invalidate(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *InMemoryCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*CacheEntry[V])
}

func (c *InMemoryCache[K, V]) Has(key K) bool {
	_, exists := c.Get(key)
	return exists
}
