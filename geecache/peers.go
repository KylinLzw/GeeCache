package geecache

import pb "GeeCache/geecachepb"

// 定义节点选择接口

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 定义节点获取数据接口

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
