# HTTPPool

![HTTPPool](https://cdn.jsdelivr.net/gh/KylinLzw/MarkdownImage/img/20230829151249.png)

- hTTPPool 中包含服务端和客户端
- 接口PeerPicker的PickPeer()方法用于根据传入的key选择相应节点PeerGetter 
- 接口PeerGetter的Get()方法用于从对应group查找缓存值
- peers，类型是一致性哈希算法的Map，用来根据具体的key选择节点。 
- httpGetters，映射远程节点与对应的httpGetter。每一个远程节点对应一个httpGetter。
- RegisterPeers()方法，将实现了PeerPicker接口的HTTPPool注入到Group中。 
- getFromPeer()方法，访问远程节点peerGetter，获取缓存值
- 则HTTPPool既具备提供HTTP服务的能力，也具备根据具体的key，创建HTTP客户端从远程节点获取缓存值的能力