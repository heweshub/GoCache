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