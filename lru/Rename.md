# 最近最久未使用淘汰算法（LRU）

- 相对于仅考虑时间的FIFO和仅考虑访问频率的LFU，LRU相对平衡。
- LRU认为，如果数据最近被访问，那么将来会被访问的概率也更高。
- 实现上用一个队列模拟时间序列，如果某条记录被访问，则移动到队尾，那么队首就是最近最少访问的数据。

### 核心数据结构

![LRU核心数据结构](https://cdn.jsdelivr.net/gh/KylinLzw/MarkdownImage/img/20230829102233.png)

- 核心数据结构是采用一个哈希表和一个双向链表来实现。
- 可以通过哈希表实现根据健值查找数据，哈希表保存健值和真实数据的映射关系。
- 双向链表保存具体数据，每个结点包含健值和真实数据。
- 如果某条记录被访问，则移动到队尾，那么队首就是最近最少访问的数据。

? 为什么使用双链表而不使用单链表？
双链表有前驱和后继节点，删除的时间复杂度为O(1)>

? 为什么双链表中还需要存储key，而不只存储value?
删除最近最少使用节点的时候，需要通过节点获取对应的key，然后再删除哈希表中的键值对。

### 核心方法

```go
func New(maxBytes int64, onEvicted func(string, Value)) *Cache{} // 创建缓存空间
func (c *Cache) Get(key string) (value Value, ok bool){} // 根据健值获取数据
func (c *Cache) Add(key string, value Value){} // 添加，修改缓存数据
func (c *Cache) RemoveOldest(){} // 淘汰结点
```



