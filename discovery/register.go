package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type ServiceInfo struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  int64  `json:"weight"`
}

type Register struct {
	Addrs       []string
	DialTimeout int

	closeCh     chan struct{}
	leaseID     clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo ServiceInfo
	srvTTL  int64
	cli     *clientv3.Client
	logger  *zap.Logger
}

func NewRegister(addrs []string, logger *zap.Logger) *Register {
	return &Register{
		Addrs:       addrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

func (r *Register) Register(srvInfo ServiceInfo, ttl int64) (chan<- struct{}, error) {
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}
	var err error
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err := r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})

	go r.keepAlive()

	return r.closeCh, nil
}

func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	r.leaseID = leaseResp.ID

	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leaseID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildPath(r.srvInfo), string(data), clientv3.WithLease(r.leaseID))
	return err
}

func (r *Register) unregister() error {
	_, err := r.cli.Revoke(context.Background(), r.leaseID)
	return err
}

func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}
		case resp := <-r.keepAliveCh:
			if resp == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		}
	}
}

func (r *Register) GetSeriveInfo() (ServiceInfo, error) {
	resp, err := r.cli.Get(context.Background(), BuildPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}
	info := ServiceInfo{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}
	return info, nil
}
