// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	// Timer for how often to recycle the expired cache items in memory (in seconds)
	DefaultEvery = 60 // 1 minute
)

// MemoryItem stores memory cache item.
type MemoryItem struct {
	val         interface{}
	createdTime time.Time
	lifespan    time.Duration
}

func (mi *MemoryItem) isExpire() bool {
	// 0 means forever
	if mi.lifespan == 0 {
		return false
	}
	return time.Since(mi.createdTime) > mi.lifespan
}

// MemoryCache is a memory cache adapter.
// Contains a RW locker for safe map storage.
type MemoryCache struct {
	sync.RWMutex
	dur   time.Duration
	items map[string]*MemoryItem
	Every int // run an expiration check Every clock time
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() Cache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

// Get returns cache from memory.
// If non-existent or expired, return nil.
func (bc *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	bc.RLock()
	defer bc.RUnlock()
	if itm, ok := bc.items[key]; ok {
		if itm.isExpire() {
			return nil, errors.New("the key is expired")
		}
		return itm.val, nil
	}
	return nil, errors.New("the key isn't exist")
}

// GetMulti gets caches from memory.
// If non-existent or expired, return nil.
func (bc *MemoryCache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	rc := make([]interface{}, len(keys))
	keysErr := make([]string, 0)

	for i, ki := range keys {
		val, err := bc.Get(context.Background(), ki)
		if err != nil {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, err.Error()))
			continue
		}
		rc[i] = val
	}

	if len(keysErr) == 0 {
		return rc, nil
	}
	return rc, errors.New(strings.Join(keysErr, "; "))
}

// Put puts cache into memory.
// If lifespan is 0, it will never overwrite this value unless restarted
func (bc *MemoryCache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items[key] = &MemoryItem{
		val:         val,
		createdTime: time.Now(),
		lifespan:    timeout,
	}
	return nil
}

// Delete cache in memory.
func (bc *MemoryCache) Delete(ctx context.Context, key string) error {
	bc.Lock()
	defer bc.Unlock()
	if _, ok := bc.items[key]; !ok {
		return errors.New("key not exist")
	}
	delete(bc.items, key)
	if _, ok := bc.items[key]; ok {
		return errors.New("delete key error")
	}
	return nil
}

// Incr increases cache counter in memory.
// Supports int,int32,int64,uint,uint32,uint64.
func (bc *MemoryCache) Incr(ctx context.Context, key string) error {
	bc.Lock()
	defer bc.Unlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch val := itm.val.(type) {
	case int:
		itm.val = val + 1
	case int32:
		itm.val = val + 1
	case int64:
		itm.val = val + 1
	case uint:
		itm.val = val + 1
	case uint32:
		itm.val = val + 1
	case uint64:
		itm.val = val + 1
	default:
		return errors.New("item val is not (u)int (u)int32 (u)int64")
	}
	return nil
}

// Decr decreases counter in memory.
func (bc *MemoryCache) Decr(ctx context.Context, key string) error {
	bc.Lock()
	defer bc.Unlock()
	itm, ok := bc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch val := itm.val.(type) {
	case int:
		itm.val = val - 1
	case int64:
		itm.val = val - 1
	case int32:
		itm.val = val - 1
	case uint:
		if val > 0 {
			itm.val = val - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint32:
		if val > 0 {
			itm.val = val - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint64:
		if val > 0 {
			itm.val = val - 1
		} else {
			return errors.New("item val is less than 0")
		}
	default:
		return errors.New("item val is not int int64 int32")
	}
	return nil
}

// IsExist checks if cache exists in memory.
func (bc *MemoryCache) IsExist(ctx context.Context, key string) (bool, error) {
	bc.RLock()
	defer bc.RUnlock()
	if v, ok := bc.items[key]; ok {
		return !v.isExpire(), nil
	}
	return false, nil
}

// ClearAll deletes all cache in memory.
func (bc *MemoryCache) ClearAll(context.Context) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items = make(map[string]*MemoryItem)
	return nil
}

// StartAndGC starts memory cache. Checks expiration in every clock time.
func (bc *MemoryCache) StartAndGC(config string) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultEvery
	}
	dur := time.Duration(cf["interval"]) * time.Second
	bc.Every = cf["interval"]
	bc.dur = dur
	go bc.vacuum()
	return nil
}

// check expiration.
func (bc *MemoryCache) vacuum() {
	bc.RLock()
	every := bc.Every
	bc.RUnlock()

	if every < 1 {
		return
	}
	for {
		<-time.After(bc.dur)
		bc.RLock()
		if bc.items == nil {
			bc.RUnlock()
			return
		}
		bc.RUnlock()
		if keys := bc.expiredKeys(); len(keys) != 0 {
			bc.clearItems(keys)
		}
	}
}

// expiredKeys returns keys list which are expired.
func (bc *MemoryCache) expiredKeys() (keys []string) {
	bc.RLock()
	defer bc.RUnlock()
	for key, itm := range bc.items {
		if itm.isExpire() {
			keys = append(keys, key)
		}
	}
	return
}

// ClearItems removes all items who's key is in keys
func (bc *MemoryCache) clearItems(keys []string) {
	bc.Lock()
	defer bc.Unlock()
	for _, key := range keys {
		delete(bc.items, key)
	}
}

func init() {
	Register("memory", NewMemoryCache)
}
