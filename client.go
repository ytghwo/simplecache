package simplecache

import (
	"context"
	"fmt"
	"time"

	"github.com/ytghwo/simplecache/registry"
	pb "github.com/ytghwo/simplecache/simplecachepb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type client struct {
	name string
}

func (c *client) Fetch(group string, key string) ([]byte, error) {
	cli, err := clientv3.New(defaultEtcdConfig)
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	conn, err := registry.EtcdDial(cli, c.name)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	grpcClient := pb.NewSimpleCacheClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := grpcClient.Get(ctx, &pb.GetRequest{
		Group: group,
		Key:   key,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get %s/%s from peer %s", group, key, c.name)
	}
	return resp.Value, nil
}

func NewClient(service string) *client {
	return &client{name: service}
}
