# 单机并发缓存 - cache

- cache用sync.Mutex封装LRU的方法，使之支持并发读写。 
- 先抽象一个只读的数据结构ByteView用来表示缓存值。
  - ByteView 只有一个数据成员，b []byte，b 将会存储真实的缓存值。
  - 并且b是只读的，返回的是b的一个拷贝，防止缓存值被外部程序修改。
- 为lru添加并发读写，cache.go实例化了lru,并封装了add()、get()方法。
  - 并且封装的方法是私有的，不会被其他包使用。
  - add方法采用延迟初始化用于提高性能，并减少程序内存要求。


# 缓存主体结构 - Group

![数据获取流程](https://cdn.jsdelivr.net/gh/KylinLzw/MarkdownImage/img/20230829105459.png)

- 一个Group可以认为是一个缓存的命名空间：
  - 第一个属性name表示每个Group拥有一个唯一的名称name。
  - 第二个属性getter Getter，即缓存未命中时获取源数据的回调。
  - 第三个属性mainCache cache，即并发缓存。
  - 第四个属性peers保存其它节点的数据获取
  - 第五个属性loader用于保证并发访问其他节点时只会发送一次http请求
- Get方法用于获取数据，可以从缓存/本地/其他节点获取到数据