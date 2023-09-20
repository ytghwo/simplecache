package simplecache

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ytghwo/simplecache/consistenthash"
	"github.com/ytghwo/simplecache/registry"
	pb "github.com/ytghwo/simplecache/simplecachepb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const (
	defaultAddr     = "127.0.0.1:6324"
	defaultReplicas = 50
)

var (
	defaultEtcdConfig = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
)

type server struct {
	pb.UnimplementedSimpleCacheServer

	addr       string
	status     bool
	stopSignal chan error
	mu         sync.Mutex
	consHash   *consistenthash.Consistency
	clients    map[string]*client
}

func NewServer(addr string) (*server, error) {
	if addr == "" {
		addr = defaultAddr
	}
	if !validPeerAddr(addr) {
		return nil, fmt.Errorf("the format must be x.x.x.x:port")
	}
	return &server{addr: addr}, nil
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	group, key := in.GetGroup(), in.GetKey()
	resp := &pb.GetResponse{}

	log.Printf("[simplecache_server %s] Recv RPC Request - (%s)/(%s)", s.addr, group, key)
	if key == "" {
		return resp, fmt.Errorf("key required")
	}
	g := GetGroup(key)
	if g == nil {
		return resp, fmt.Errorf("group not found")
	}
	view, err := g.Get(key)
	if err != nil {
		return resp, err
	}
	resp.Value = view.ByteSlice()
	return resp, nil
}

func (s *server) Start() error {
	s.mu.Lock()
	if s.status {
		s.mu.Unlock()
		return fmt.Errorf("server already started")
	}
	s.status = true
	s.stopSignal = make(chan error)
	port := strings.Split(s.addr, ":")[1]
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleCacheServer(grpcServer, s)
	go func() {
		err := registry.Register("simplecache", s.addr, s.stopSignal)
		if err != nil {
			log.Fatalf(err.Error())
		}
		close(s.stopSignal)
		err = lis.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("[%s] Revoke service and close tcp socket ok.", s.addr)
	}()
	s.mu.Unlock()
	if err := grpcServer.Serve(lis); s.status && err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func (s *server) Pick(key string) (Fetcher, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	peerAddr := s.consHash.GetPeer(key)
	if peerAddr == s.addr {
		log.Printf("ooh! pick myself,I am %s\n", s.addr)
		return nil, false
	}
	log.Printf("[cache %s] pick remote peer: %s\n", s.addr, peerAddr)
	return s.clients[peerAddr], true
}

func (s *server) SetPeers(peerAddr ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consHash = consistenthash.New(defaultReplicas, nil)
	s.consHash.Register(peerAddr...)
	s.clients = make(map[string]*client)
	for _, peerAddr := range peerAddr {
		if !validPeerAddr(peerAddr) {
			panic(fmt.Sprintf("[peer %s] invalid address format,int should be x.x.x.x:port", peerAddr))
		}
		service := fmt.Sprintf("simplecache/%s", peerAddr)
		s.clients[peerAddr] = NewClient(service)
	}
}

func (s *server) Stop() {
	s.mu.Lock()
	if !s.status {
		s.mu.Unlock()
		return
	}
	s.stopSignal <- nil
	s.status = false
	s.clients = nil
	s.consHash = nil
	s.mu.Unlock()
}
