### 极客兔兔 7天用Go从零实现分布式缓存GeeCache

### https://geektutu.com/post/geecache.html

#### 跟着教程动手实践
### 总结
- 7 天用 Go 动手写/从零实现分布式缓存 GeeCache 这个系列就完成了。简单回顾下。
- 第一天，为了解决资源限制的问题，实现了 LRU 缓存淘汰算法；
- 第二天实现了单机并发，并给用户提供了自定义数据源的回调函数；
- 第三天实现了 HTTP 服务端；
- 第四天实现了一致性哈希算法，解决远程节点的挑选问题；
- 第五天创建 HTTP 客户端，实现了多节点间的通信；
- 第六天实现了 singleflight 解决缓存击穿的问题；
- 第七天，使用 protobuf 库，优化了节点间通信的性能。