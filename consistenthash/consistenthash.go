package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 自定义哈希函数类型
type Hash func(data []byte) uint32

// Map 一致性哈希数据结构
type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点扩充倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 哈希值和节点的对应关系
}

// New 新建一个一致性哈希数据结构
func New(replicas int, hf Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hf,
		hashMap:  make(map[int]string),
	}

	// 默认采用 crc32.ChecksumIEEE
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add  添加结点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 扩充虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	// 哈希环排序
	sort.Ints(m.keys)
}

// Get 获取缓存值存储的结点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	// 二分查找服务的节点
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 最后部分需要第一个节点服务，进行取余
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

// Remove 删除节点
func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
}
