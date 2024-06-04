// https://github.com/bytedance/gopkg/blob/develop/util/gopool

package gopool

import (
	"context"
	"fmt"
	"sync"
)

var defaultPool Pool

var poolMap sync.Map

func init() {
	defaultPool = NewPool("gopool.DeafultPool", 50, NewConfig())
}

func Go(f func()) {
	CtxGo(context.Background(), f)
}

func CtxGo(ctx context.Context, f func()) {
	defaultPool.CtxGo(ctx, f)
}

func SetCap(cap int32) {
	defaultPool.SetCap(cap)
}

func SetPanicHandler(f func(context.Context, any)) {
	defaultPool.SetPanicHandler(f)
}

func WorkerCount() int32 {
	return defaultPool.WorkCount()
}

func RegisterPool(p Pool) error {
	_, loaded := poolMap.LoadOrStore(p.Name(), p)
	if loaded {
		return fmt.Errorf("name: %s already registered", p.Name())
	}
	return nil
}

func GetPool(name string) Pool {
	p, ok := poolMap.Load(name)
	if !ok {
		return nil
	}
	return p.(Pool)
}
