/*
极客兔兔 :动手写分布式缓存 - GeeCache第一天 LRU 缓存淘汰
https://geektutu.com/post/geecache-day1.html
*/
package lru

import (
	"container/list"
)

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 获取元素
func (c *Cache) Get(key string) (value Value, ok bool) {
	var (
		ele *list.Element
	)
	if ele, ok = c.cache[key]; ok {
		c.ll.MoveToBack(ele) // 最近访问的元素移动到队尾
		kv := ele.Value.(*entry)
		return kv.value, ok
	}
	return
}

// 删除最老的元素 即队列头的元素
func (c *Cache) removeOldest() {
	var (
		ele *list.Element
	)
	if ele = c.ll.Front(); ele != nil {
		c.ll.Remove(ele) // 双向列表删除元素
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // cache删除元素
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 修改已缓存数据大小
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) // 执行删除元素的回调
		}
	}
}

// 添加元素 添加到队列尾部 同时需要判断是否需要移除元素
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 如果key存在则把元素移动到队列尾部
		c.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		// 更新已缓存数据大小
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 如果不存在则在队尾插入元素
		ele = c.ll.PushBack(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 判断是否需要移除元素
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.removeOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
