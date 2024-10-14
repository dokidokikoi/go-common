package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/dokidokikoi/go-common/discovery/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var etcdAddrs = []string{
	"127.0.0.1:23791",
	"127.0.0.1:23792",
	"127.0.0.1:23793",
}

func Testesolver(t *testing.T) {
	r := NewResolver(etcdAddrs)
	resolver.Register(r)

	// etcd中注册5个服务
	go newServer(t, 1001, "1.0.0", 1)
	go newServer(t, 1002, "1.0.0", 1)
	go newServer(t, 1003, "1.0.0", 1)
	go newServer(t, 1004, "1.0.0", 1)
	go newServer(t, 1006, "1.0.0", 10)

	conn, err := grpc.Dial("etcd:///hello", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		t.Fatalf("failed to dial %v", err)
	}
	defer conn.Close()

	c := testdata.NewGreeterClient(conn)

	// 进行十次数据请求
	for i := 0; i < 10; i++ {
		resp, err := c.SayHello(context.Background(), &testdata.HelloRequest{Name: "abc"})
		if err != nil {
			t.Fatalf("say hello failed %v", err)
		}
		log.Println(resp.Message)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)
}

type server struct {
	Port int
	testdata.UnimplementedGreeterServer
}

func (s *server) SayHello(c context.Context, req *testdata.HelloRequest) (*testdata.HelloReply, error) {
	return &testdata.HelloReply{Message: fmt.Sprintf("Hello from %d", s.Port)}, nil
}

func newServer(t *testing.T, port int, version string, weight int64) {
	register := NewRegister(etcdAddrs)
	defer register.Stop()

	listen, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	testdata.RegisterGreeterServer(s, &server{Port: port})

	info := ServiceInfo{
		Name:    "test",
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Version: version,
		Weight:  weight,
	}

	register.Register(info, 10)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
