package simplecache

import (
	"fmt"
	"log"
	"sync"

	"github.com/ytghwo/simplecache/singleflight"
)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Retriever interface {
	retrieve(string) ([]byte, error)
}

type RetrieverFunc func(key string) ([]byte, error)

func (f RetrieverFunc) retrieve(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	cache     *cache
	retriever Retriever
	server    Picker
	flight    *singleflight.Flight
}

func NewGroup(name string, maxBytes int64, retriever Retriever) *Group {
	if retriever == nil {
		panic("Group retriever must be existed!")
	}
	g := &Group{
		name:      name,
		cache:     NewCache(int(maxBytes)),
		retriever: retriever,
		flight:    &singleflight.Flight{},
	}
	mu.Lock()
	groups[name] = g
	mu.Unlock()
	return g
}

func (g *Group) RegisterSvr(p Picker) {
	if g.server != nil {
		panic("group had benn registered server")
	}
	g.server = p
}

func GetGroup(name string) *Group {
	mu.Lock()
	g := groups[name]
	mu.Unlock()
	return g
}

func DestoryGroup(name string) {
	g := GetGroup(name)
	if g != nil {
		svr := g.server.(*server)
		svr.Stop()
		delete(groups, name)
		log.Printf("Destory cache [%s %s]", name, svr.addr)
	}
}

func (g *Group) Get(key string) (byteview, error) {
	if key == "" {
		return byteview{}, fmt.Errorf("key required")
	}
	if value, ok := g.cache.get(key); ok {
		log.Println("cache hit")
		return value, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (byteview, error) {
	view, err := g.flight.Fly(key, func() (interface{}, error) {
		if g.server != nil {
			if fetcher, ok := g.server.Pick(key); ok {
				bytes, err := fetcher.Fetch(g.name, key)
				if err == nil {
					return byteview{b: cloneBytes(bytes)}, nil
				}
				log.Printf("fail to get *%s from peer,%s.\n", key, err.Error())
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return view.(byteview), err
	}
	return byteview{}, err
}

func (g *Group) getLocally(key string) (byteview, error) {
	bytes, err := g.retriever.retrieve(key)
	if err != nil {
		return byteview{}, err
	}
	value := byteview{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value byteview) {
	g.cache.add(key, value)
}
