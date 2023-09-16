package consistenthash

import (
	"hash/crc32"
	"log"
	"sort"
	"testing"
)

func TestRegister(t *testing.T) {
	c := New(2, nil)
	c.Register("peer1", "peer2")
	if len(c.ring) != 5 {
		t.Errorf("Actual: %d Expect: %d", len(c.ring), 5)
	}
	hashValue := int(crc32.ChecksumIEEE([]byte("1peer1")))
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= hashValue
	})
	if c.ring[idx] != hashValue {
		t.Errorf("Actual: %d Expect: %d", c.ring[idx], hashValue)
	}
}

func TestGet(t *testing.T) {
	c := New(2, nil)
	c.Register("peer1", "peer2")
	key := "tom"
	hashValue := int(crc32.ChecksumIEEE([]byte(key)))
	log.Printf("key hashvalue: %d\n", hashValue)
	for _, v := range c.ring {
		log.Printf("%d -> %s\n", v, c.hashMap[v])
	}
	peer := c.GetPeer(key)
	log.Printf("go to search %s", peer)
}
