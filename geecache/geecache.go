package geecache

import (
	pb "GeeCache/geecachepb"
	"GeeCache/singleflight"
	"errors"
	"fmt"
	"log"
	"sync"
)

// Getter 函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，
// 也能够传入实现了该接口的结构体作为参数。
// 接口型函数的接口类型只能有一个方法，这样才能让某一个函树类型实现该接口。
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 实现 Getter 接口的 Get 方法
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group ：每一块缓存区
// 第一个属性name表示每个Group拥有一个唯一的名称name。
// 第二个属性getter Getter，即缓存未命中时获取源数据的回调。
// 第三个属性mainCache cache，即并发缓存。
// 第四个属性peers用于作为客户端获取其他节点
// 第五个属性loader用于保证并发访问其他节点时只会发送一次http请求
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	hotCache  cache
	cacheByte int64
	peers     PeerPicker
	loader    *singleflight.Group
}

// 访问缓存区的全局读写锁和缓存区名称对应关系
// 加锁保证并发互斥访问
var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// GetGroup 根据缓存区名称获取具体的缓存区
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	group := groups[name]
	return group
}

// NewGroup 申请新的一块缓存区
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		cacheByte: cacheBytes,
		mainCache: cache{cacheByte: cacheBytes},
		hotCache:  cache{cacheByte: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// Get : 获取数据
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 主缓存命中
	if v, ok := g.mainCache.get(key); ok {
		log.Println("Main GeeCache hit...")
		return v, nil
	}

	// 热点缓存命中
	if v, ok := g.hotCache.get(key); ok {
		log.Println("Hot GeeCache hit...")
		return v, nil
	}

	// 缓存未命中
	return g.load(key)
}

// getLocally : 本地结点获取数据
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)

	// 数据库查找不到，在缓存中加入标记点
	if err == errors.New("key Not Exist") {
		bytes = []byte("key Not Exist")
	}
	if err != nil && err != errors.New("key Not Exist") {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value, &g.mainCache)
	return value, nil
}

// populateCache ：数据加入到缓存中
func (g *Group) populateCache(key string, value ByteView, cache *cache) {
	if g.cacheByte <= 0 {
		return
	}
	cache.add(key, value)

	for {
		mainBytes := g.mainCache.bytes()
		hotBytes := g.hotCache.bytes()
		if mainBytes+hotBytes <= g.cacheByte {
			return
		}

		victim := &g.mainCache
		if hotBytes > mainBytes/8 {
			victim = &g.hotCache
		}
		victim.removeOldest()
	}
}

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// 从其他节点获取数据
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)

	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

// 缓存未命中时获取数据
func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					g.populateCache(key, value, &g.hotCache)
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}
