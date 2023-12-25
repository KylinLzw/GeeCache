package geecache

// 抽象一个只读数据结构 ByteView 用来表示缓存值
// ByteView 只有一个数据成员，b []byte，b 将会存储真实的缓存值。
// 实现 Len() int 方法，我们在 lru.Cache 的实现中，要求被缓存对象必须实现 Value 接口，
// b 是只读的，使用 cloneBytes() 方法返回一个拷贝，防止缓存值被外部程序修改。

type ByteView struct {
	b []byte
}

func (B ByteView) Len() int {
	return len(B.b)
}

func (B ByteView) String() string {
	return string(B.b)
}

func (B ByteView) ByteSlice() []byte {
	return cloneBytes(B.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
