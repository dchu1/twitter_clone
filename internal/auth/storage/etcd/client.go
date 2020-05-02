package etcd

import (
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var Client *clientv3.Client

func NewClient(endpoints []string) (*clientv3.Client, error) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli, err
}
