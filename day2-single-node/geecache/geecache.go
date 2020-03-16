/*
极客兔兔
动手写分布式缓存 - GeeCache第二天 单机并发缓存
https://geektutu.com/post/geecache-day2.html
*/
package geecache

import (
	"fmt"
	"log"
	"sync"
)

var (
	mu     sync.Mutex
	groups = make(map[string]*Group) // 可以理解为一个数据库
)

// 可以理解为一个表
type Group struct {
	name      string
	getter    Getter // 如果没有缓存则调用getter获取原始值 比如从数据库获取
	mainCache cache
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
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
func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := groups[name]
	return g
}

//                            是
//接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
//                |  否                         是
//                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
//                            |  否
//                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit ", key)
		return v, nil
	}
	return g.load(key)
}
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	// 将key-value放入缓存
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
