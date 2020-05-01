package etcd

import (
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var Client *clientv3.Client

func NewClient(endpoints []string) (*clientv3.Client, error) {
	// // expect dial time-out on ipv4 blackhole
	// _, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{"http://254.0.0.1:12345"},
	// 	DialTimeout: 2 * time.Second,
	// })

	// // etcd clientv3 >= v3.2.10, grpc/grpc-go >= v1.7.3
	// if err == context.DeadlineExceeded {
	// 	log.Fatal(err)
	// }

	// // etcd clientv3 <= v3.2.9, grpc/grpc-go <= v1.2.1
	// if err == grpc.ErrClientConnTimeout {
	// 	log.Fatal(err)
	// }

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli, err
}
