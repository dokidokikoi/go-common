package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

func BuildPrefix(info ServiceInfo) string {
	if info.Version == "" {
		return fmt.Sprintf("/%s/", info.Name)
	}
	return fmt.Sprintf("/%s/%s", info.Name, info.Version)
}

func BuildPath(info ServiceInfo) string {
	return fmt.Sprintf("%s/%s", BuildPrefix(info), info.Addr)
}

func ParseValue(value []byte) (ServiceInfo, error) {
	var info ServiceInfo
	err := json.Unmarshal(value, &info)
	return info, err
}

func SplitPath(path string) (ServiceInfo, error) {
	var info ServiceInfo
	strs := strings.Split(path, "/")
	if len(strs) < 1 {
		return info, errors.New("invalid path")
	}
	info.Addr = strs[len(strs)-1]
	return info, nil
}

func Exist(addrs []resolver.Address, addr resolver.Address) bool {
	for i := range addrs {
		if addrs[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

func Remove(addrs []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range addrs {
		if addrs[i].Addr == addr.Addr {
			addrs[i] = addrs[len(addrs)-1]
			return addrs[:len(addrs)-1], true
		}
	}
	return nil, false
}

// func BuildResolverUrl(app string) string {
// 	return sc
// }
