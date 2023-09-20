package lru

import (
	"container/list"
)

//lru模块，实现缓存淘汰机制

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

type onEliminated func(key string, value Value)

type Cache struct {
	capacity int64
	length   int64
	list     *list.List
	hash     map[string]*list.Element
	callBack onEliminated
}

func New(capacity int64, callback onEliminated) *Cache {
	return &Cache{
		capacity: capacity,
		list:     list.New(),
		hash:     make(map[string]*list.Element),
		callBack: callback,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.hash[key]; ok {
		c.list.MoveToFront(elem)
		entry := elem.Value.(*entry)
		return entry.value, true
	}
	return
}

func (c *Cache) Add(key string, value Value) {
	kvsize := int64(len(key)) + int64(value.Len())
	if c.capacity != 0 && kvsize+int64(c.length) > int64(c.capacity) {
		c.remove()
	}
	if elem, ok := c.hash[key]; ok {
		c.list.MoveToFront(elem)
		oldEntry := elem.Value.(*entry)
		c.length += int64(value.Len()) - int64(oldEntry.value.Len())
		oldEntry.value = value
	} else {
		elem := c.list.PushFront(&entry{key: key, value: value})
		c.hash[key] = elem
		c.length += kvsize
	}
}

func (c *Cache) remove() {
	elem := c.list.Back()
	if elem != nil {
		entry := elem.Value.(*entry)
		k, v := entry.key, entry.value
		delete(c.hash, k)
		c.list.Remove(elem)
		c.length -= int64(len(k)) + int64(v.Len())
		if c.callBack != nil {
			c.callBack(k, v)
		}
	}
}
