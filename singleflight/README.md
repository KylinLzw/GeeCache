# SingleFlight

### 解决问题：
- 如果并发了N个相同请求（注意是同时发送请求）。 
- 假设对数据库的访问没有做任何限制，很有可能向数据库也发起N次请求，容易导致缓存击穿或穿透。
- 即使对数据库做了防护，HTTP请求也是非常消耗资源的操作，针对相同的key,也没有必要向远程缓存节点发起N次相同请求。
- 针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误。