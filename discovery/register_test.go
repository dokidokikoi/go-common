package discovery

import (
	"fmt"
	"testing"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func TestRegister(t *testing.T) {
	info := ServiceInfo{
		Name:    "user",
		Addr:    "localhost:8083",
		Version: "1.0.0",
		Weight:  2,
	}
	addrs := []string{
		"127.0.0.1:23791",
		"127.0.0.1:23792",
		"127.0.0.1:23793",
	}
	r := NewRegister(addrs, zaplog.L())

	_, err := r.Register(info, 120)
	if err != nil {
		t.Fatalf("register to etcd failed %v", err)
	}

	infoRes, err := r.GetSeriveInfo()
	if err != nil {
		t.Fatalf("get info failed %v", err)
	}
	fmt.Printf("%v\n12", infoRes)
}
