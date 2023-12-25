package geecache

import (
	"GeeCache/lru"
	"sync"
)

// cache.go 的实现非常简单，实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu。

type cache struct {
	mu        sync.Mutex
	lru       *lru.Cache
	cacheByte int64
}

// 在 add 方法中，判断了 c.lru 是否为 nil，如果等于 nil 再创建实例。这种方法称之为延迟初始化(Lazy Initialization)，
// 一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。主要用于提高性能，并减少程序内存要求。
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheByte, nil)
	}
	c.lru.Add(key, value)
}

// 通过加锁，保证对数据的并发访问
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

// 获取现在缓存区存储数据的大小
func (c *cache) bytes() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return 0
	}

	return c.lru.Bytes()
}

// 删除缓存数据
func (c *cache) removeOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	c.lru.RemoveOldest()
}
