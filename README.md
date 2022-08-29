# GoCache
disturbuted cache system based on Golang

分布式缓存  
> 很多高并发的操作都是需要借助cache+MQ来实现的
> 缓存淘汰策略
> 并发伴随着读写冲突
> 单机性能瓶颈
    可以分为水平扩展和垂直扩展


GoCache
1. 单机缓存和基于HTTP的分布式缓存
2. LRU缓存策略
3. Go锁机制防止缓存击穿
4. 一致性哈希选择节点，实现负载均衡
5. 使用pb优化节点间二进制通信


实现策略
1. LRU
    核心数据结构 map + 双向链表

2. sync.Mutex 实现LRU缓存的并发控制
    实现GoCache的核心数据结构Group 

3. 搭建HTTP server
    http.ListenAndServe 接收 2 个参数，第一个参数是服务启动的地址，第二个参数是 Handler，任何实现了 ServeHTTP 方法的对象都可以作为 HTTP 的 Handler。
    通过URL的路径来check搜索的key有没有在cache中（cache用的是LRU算法实现的）

4. 一致性哈希
    通过哈希每次访问同一个节点
    使用一致性哈希算法可以解决存在的缓存雪崩
    引入环的概念，解决了增加或者删除节点之后节点挂载问题
    对于数据倾斜问题，提出了虚拟节点的概念

5. 分布式节点
    注册节点，借助一致性哈希算法选择节点
    使用HTTP与远程节点实现通信(从远节点获取缓存值)

6. 解决缓存击穿问题
    singleflight

7. PB通信
    protobuf 即 Protocol Buffers，Google 开发的一种数据描述语言，是一种轻便高效的结构化数据存储格式，与语言、平台无关，可扩展可序列化。protobuf 以二进制方式存储，占用空间小。
    RPC用的会比较多
