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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DiskCache[K comparable, V any] struct {
	mu        sync.RWMutex
	cacheDir  string
	namespace string
	items     map[K]*CacheEntry[V]
	dirty     bool
}

func NewDiskCache[K comparable, V any](namespace string) *DiskCache[K, V] {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// fmt.Printf("failed to get home directory: %v", err)
		return nil
	}

	cacheDir := filepath.Join(homeDir, ".cache", "catalyst")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		// fmt.Printf("failed to create cache directory: %v", err)
		return nil
	}

	cache := &DiskCache[K, V]{
		cacheDir:  cacheDir,
		namespace: namespace,
		items:     make(map[K]*CacheEntry[V]),
	}

	if err := cache.load(); err != nil {
		cache.items = make(map[K]*CacheEntry[V])
	}

	return cache
}

func (c *DiskCache[K, V]) getCacheFilePath() string {
	return filepath.Join(c.cacheDir, c.namespace+".json")
}

func (c *DiskCache[K, V]) load() error {
	cachePath := c.getCacheFilePath()

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	if err := json.Unmarshal(data, &c.items); err != nil {
		return fmt.Errorf("failed to parse cache file: %w", err)
	}

	return nil
}

func (c *DiskCache[K, V]) save() error {
	cachePath := c.getCacheFilePath()

	data, err := json.MarshalIndent(c.items, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	c.dirty = false
	return nil
}

func (c *DiskCache[K, V]) Get(key K) (V, bool) {
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

func (c *DiskCache[K, V]) Set(key K, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheEntry[V]{
		Data:     value,
		CachedAt: time.Now(),
	}
	c.dirty = true

	return c.save()
}

func (c *DiskCache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	c.items[key] = &CacheEntry[V]{
		Data:      value,
		CachedAt:  time.Now(),
		ExpiresAt: &expiresAt,
	}
	c.dirty = true

	return c.save()
}

func (c *DiskCache[K, V]) Invalidate(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	c.dirty = true

	return c.save()
}

func (c *DiskCache[K, V]) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*CacheEntry[V])
	c.dirty = true

	return c.save()
}

func (c *DiskCache[K, V]) Has(key K) bool {
	_, exists := c.Get(key)
	return exists
}

// This is useful when we need to read, modify, and save in one atomic operation
func (c *DiskCache[K, V]) Update(key K, updateFn func(V) V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.items[key]
	var currentValue V
	if exists {
		currentValue = entry.Data
	}

	newValue := updateFn(currentValue)

	c.items[key] = &CacheEntry[V]{
		Data:     newValue,
		CachedAt: time.Now(),
	}
	c.dirty = true

	return c.save()
}
