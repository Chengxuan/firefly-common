// Copyright © 2023 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
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

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCacheReturnsSameCacheForSameConfig(t *testing.T) {
	ctx := context.Background()
	cacheManager := NewCacheManager(ctx, true)
	cache0, _ := cacheManager.GetCache(ctx, "cacheA", 85, time.Second, true)
	cache1, _ := cacheManager.GetCache(ctx, "cacheA", 85, time.Second, true)

	assert.Equal(t, cache0, cache1)
	assert.Equal(t, []string{"cacheA"}, cacheManager.ListCacheNames())

	cache2, _ := cacheManager.GetCache(ctx, "cacheB", 85, time.Second, true)
	assert.NotEqual(t, cache0, cache2)
	assert.Equal(t, 2, len(cacheManager.ListCacheNames()))
}

func TestTwoSeparateCacheWorksIndependently(t *testing.T) {
	ctx := context.Background()
	cacheManager := NewCacheManager(ctx, true)
	cache0, _ := cacheManager.GetCache(ctx, "cacheA", 85, time.Second, true)
	cache1, _ := cacheManager.GetCache(ctx, "cacheB", 85, time.Second, true)

	cache0.SetInt("int0", 100)
	assert.Equal(t, 100, cache0.GetInt("int0"))
	assert.Equal(t, 100, cache0.Get("int0").(int))
	assert.Equal(t, 0, cache1.GetInt("int0"))
	assert.Equal(t, nil, cache1.Get("int0"))

	cache1.SetString("string1", "val1")
	assert.Equal(t, "", cache0.GetString("string1"))
	assert.Equal(t, nil, cache0.Get("string1"))
	assert.Equal(t, "val1", cache1.GetString("string1"))
	assert.Equal(t, "val1", cache1.Get("string1").(string))
	deleted := cache1.Delete("string1")
	assert.True(t, deleted)
	assert.Equal(t, nil, cache1.Get("string1"))
}

func TestCacheEnablement(t *testing.T) {
	ctx := context.Background()
	var hundred int64 = 100
	var zero int64 = 0
	// test global enablement set to false
	disabledCacheManager := NewCacheManager(ctx, false)
	cache0, _ := disabledCacheManager.GetCache(ctx, "cache0", 85, time.Second, false)
	assert.Equal(t, false, cache0.IsEnabled())

	cache0.SetInt64("int0", hundred)
	assert.Equal(t, nil, cache0.Get("int0"))

	// check individual cache cannot be turned on when the cache manager is disabled
	cache1, _ := disabledCacheManager.GetCache(ctx, "cache1", 85, time.Second, true)
	assert.Equal(t, false, cache1.IsEnabled())

	// test global enablement set to true

	enabledCacheManager := NewCacheManager(ctx, true)
	// check individual cache can be turned off when the cache manager is enabled
	cache0, _ = enabledCacheManager.GetCache(ctx, "cache0", 85, time.Second, false)
	assert.Equal(t, false, cache0.IsEnabled())

	cache0.SetInt64("int0", hundred)
	assert.Equal(t, zero, cache0.GetInt64("int0"))
	deleted := cache0.Delete("int0")
	assert.False(t, deleted)

	cache1, _ = enabledCacheManager.GetCache(ctx, "cache1", 85, time.Second, true)
	assert.Equal(t, true, cache1.IsEnabled())

	cache1.SetInt64("int0", hundred)
	assert.Equal(t, hundred, cache1.GetInt64("int0"))

}

func TestUmmanagedCacheInstance(t *testing.T) {
	uc0 := NewUmanagedCache(context.Background(), 100, 5*time.Minute)
	uc1 := NewUmanagedCache(context.Background(), 100, 5*time.Minute)
	assert.NotEqual(t, uc0, uc1)
}