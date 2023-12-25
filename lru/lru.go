package lru

import "container/list"

// Value : 实际需求的值
// 为了通用性，我们允许值是实现了 Value 接口的任意类型
// 该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// entry : 链表结点
// 键值对 entry 是双向链表节点的数据类型
// 保存key是为了便于删除反射查找相对应的map
type entry struct {
	key   string
	value Value
}

// Cache : 缓存空间
// maxBytes 是允许使用的最大内存，nBytes 是当前已使用的内存，
// 直接使用 Go 语言标准库实现的双向链表list.List。
// 字典map[string]*list.Element，键是字符串，值是双向链表中对应节点的指针
// OnEvicted 是某条记录被移除时的回调函数，可以为 nil。
type Cache struct {
	maxBytes int64
	nBytes   int64

	ll    *list.List
	cache map[string]*list.Element

	OnEvicted func(key string, value Value)
}

// New : 创建缓存
// 赋予map和list内存空间，初始化最大空间
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get ：查找功能
// 1. 通过map查找双向链表对应的结点
// 2. 将对应结点移动到队尾
// 3. 放回查找是否成功和相对应的结点value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, ok
	}

	return
}

// RemoveOldest : 淘汰最近最少访问的结点
// 也就是删除链表的尾结点，更新已经使用的空间，删除对应的map结点
// 如果回调函数 OnEvicted 不为 nil，则调用回调函数。
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len : 获取存储了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}

// Add ： 添加、修改数据
// 如果键存在，则更新对应节点的值，并将该节点移到队尾。
// 不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
// 更新 c.nBytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) + int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Bytes() int64 {
	return c.nBytes
}
