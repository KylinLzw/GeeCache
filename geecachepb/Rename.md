# protobuf

- protobuf 广泛地应用于远程过程调用(RPC) 的二进制传输，使用 protobuf 的目的非常简单，为了获得更高的性能。
- 传输前使用 protobuf 编码，接收方再进行解码，可以显著地降低二进制传输的大小。
- 另外一方面，protobuf 可非常适合传输结构化数据，便于通信字段的扩展。
- serveHTTP()使用proto.Marshal()编码HTTP响应，Get()中使用proto.Unmarshal()解码HTTP响应。