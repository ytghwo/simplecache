package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type hashFunc func(data []byte) uint32

type Consistency struct {
	hash     hashFunc
	replicas int
	ring     []int
	hashMap  map[int]string
}

func New(replicas int, fn hashFunc) *Consistency {
	c := &Consistency{
		hash:     fn,
		replicas: replicas,
		ring:     make([]int, 0),
		hashMap:  make(map[int]string),
	}
	if fn == nil {
		c.hash = crc32.ChecksumIEEE
	}
	return c
}

func (c *Consistency) Register(peerName ...string) {
	for _, peerName := range peerName {
		for i := 0; i < c.replicas; i++ {
			hashValue := int(c.hash([]byte(strconv.Itoa(i) + peerName)))
			c.ring = append(c.ring, hashValue)
			c.hashMap[hashValue] = peerName
		}
	}
	sort.Ints(c.ring)
}

func (c *Consistency) GetPeer(key string) string {
	if len(c.ring) == 0 {
		return ""
	}
	hashValue := int(c.hash([]byte(key)))
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= hashValue
	})
	return c.hashMap[c.ring[idx%len(c.ring)]]
}
