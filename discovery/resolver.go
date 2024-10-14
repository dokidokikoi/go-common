package discovery

import (
	"context"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

const (
	scheme = "etcd"
)

var _ resolver.Resolver = (*Resolver)(nil)
var _ resolver.Builder = (*Resolver)(nil)

type Resolver struct {
	scheme      string
	Addrs       []string
	DialTimeout int

	closeCh   chan struct{}
	watchCh   clientv3.WatchChan
	cli       *clientv3.Client
	keyPrefix string
	srvAddrs  []resolver.Address

	cc resolver.ClientConn
}

func NewResolver(addrs []string) *Resolver {
	return &Resolver{
		scheme:      scheme,
		Addrs:       addrs,
		DialTimeout: 3,
	}
}

func (r *Resolver) Scheme() string {
	return r.scheme
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	r.keyPrefix = BuildPrefix(ServiceInfo{Name: target.Endpoint()})
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

func (r *Resolver) start() (chan<- struct{}, error) {
	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	resolver.Register(r)

	r.closeCh = make(chan struct{})

	if err = r.sync(); err != nil {
		return nil, err
	}

	go r.watch()

	return r.closeCh, nil
}

func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Minute)
	r.watchCh = r.cli.Watch(context.Background(), r.keyPrefix, clientv3.WithPrefix())

	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {

			}
		}
	}
}

func (r *Resolver) update(events []*clientv3.Event) {
	for _, e := range events {
		var info ServiceInfo
		var err error

		switch e.Type {
		case mvccpb.PUT:
			info, err = ParseValue(e.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
			if !Exist(r.srvAddrs, addr) {
				r.srvAddrs = append(r.srvAddrs, addr)
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrs})
			}
		case mvccpb.DELETE:
			info, err = ParseValue(e.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr}
			if s, ok := Remove(r.srvAddrs, addr); ok {
				r.srvAddrs = s
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrs})
			}
		}
	}
}

func (r *Resolver) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.cli.Get(ctx, r.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	r.srvAddrs = []resolver.Address{}

	for _, v := range res.Kvs {
		info, err := ParseValue(v.Value)
		if err != nil {
			continue
		}
		addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
		r.srvAddrs = append(r.srvAddrs, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: r.srvAddrs})
	return nil
}
