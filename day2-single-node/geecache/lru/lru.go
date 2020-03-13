/*
@Author : hrbc
@Time : 2020/3/13 3:36 PM
*/
package lru

import "container/list"

type Cache struct {
	maxBytes int64      // 最大容量
	nBytes   int64      // 当前容量
	ll       *list.List // 双向队列 Back表示最新一端
	cache    map[string]*list.Element
	OnEvict  func(key string, value Value) // 淘汰元素回调
}

func NewCache(maxBytes int64, onEvict func(key string, value Value)) *Cache {
	return &Cache{maxBytes: maxBytes, ll: list.New(), cache: make(map[string]*list.Element), OnEvict: onEvict}
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.ll.MoveToBack(ele)
		return kv.value, true
	}
	return
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele = c.ll.PushBack(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.removeOldest()
	}
}

func (c *Cache) removeOldest() {
	if ele := c.ll.Front(); ele != nil {
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.ll.Remove(ele)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvict != nil {
			c.OnEvict(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
