package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

var defaultHashFunc = crc32.ChecksumIEEE

type Hash func(data []byte) uint32

// 维护一个多节点的hash环
type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // hash环 升序
	hashMap  map[int]string // 虚拟节点与真实节点映射表 键是虚拟节点hash值 值是真实节点的名称
}

func NewMap(replicas int, hash Hash) *Map {
	m := &Map{hash: hash, replicas: replicas, hashMap: make(map[int]string)}
	if m.hash == nil {
		m.hash = defaultHashFunc
	}
	return m
}

// 添加真实节点/机器
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := m.hash([]byte(strconv.Itoa(i) + key))
			m.keys = append(m.keys, int(hash))
			m.hashMap[int(hash)] = key
		}
	}
	sort.Ints(m.keys)
}

// 根据key选择机器/节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
