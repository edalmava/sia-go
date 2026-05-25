package cache

import (
	"fmt"
	"sync"
	"time"
)

type item[T any] struct {
	value     T
	expiresAt time.Time
}

func (i *item[T]) expired() bool {
	return !i.expiresAt.IsZero() && time.Now().After(i.expiresAt)
}

type Cache[T any] struct {
	name string
	data sync.Map
	ttl  time.Duration
}

func New[T any](name string, ttl time.Duration) *Cache[T] {
	return &Cache[T]{name: name, ttl: ttl}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	raw, ok := c.data.Load(key)
	if !ok {
		var zero T
		return zero, false
	}
	entry := raw.(*item[T])
	if entry.expired() {
		c.data.Delete(key)
		var zero T
		return zero, false
	}
	return entry.value, true
}

func (c *Cache[T]) Set(key string, value T) {
	var exp time.Time
	if c.ttl > 0 {
		exp = time.Now().Add(c.ttl)
	}
	c.data.Store(key, &item[T]{value: value, expiresAt: exp})
}

func (c *Cache[T]) Delete(key string) {
	c.data.Delete(key)
}

func (c *Cache[T]) Clear() {
	c.data.Range(func(k, _ any) bool {
		c.data.Delete(k)
		return true
	})
}

func (c *Cache[T]) Key(parts ...any) string {
	if len(parts) == 0 {
		return c.name
	}
	key := c.name + ":"
	for i, p := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", p)
	}
	return key
}
