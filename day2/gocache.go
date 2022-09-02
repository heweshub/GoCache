package daytwo

import (
	"fmt"
	"log"
	"sync"
)

// 函数类型实现某一个接口，称之为接口型函数，
// 方便使用者在调用时既能够传入函数作为参数，
// 也能够传入实现了该接口的结构体作为参数。
type Getter interface { // 接口
	Get(key string) ([]byte, error)
}

// 实现接口的函数
type GetterFunc func(key string) ([]byte, error)

// 实现接口
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取源数据的回调方法callback
	mainCache cache  // 实现并发存储
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) // 全局的map
)

// 创建一个group的实例，并加入全局的groups中
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		// 缓存没有命中需要从文件或者数据中获取数据并加入到缓存中
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// 全局groups中取group
func GetGroup(name string) *Group {
	mu.RLock() // 读锁
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GoCache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// 分布式场景下回调g.getter.Get()获取源数据，并加入mainCache中
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	// 回调函数没获取到数据，就直接返回空和err
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	// g.populateCache(key, value)
	g.mainCache.add(key, value)
	return value, nil
}

// func (g *Group) populateCache(key string, value ByteView) {
// 	g.mainCache.add(key, value)
// }
