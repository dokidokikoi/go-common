package gopool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool interface {
	Name() string
	SetCap(cap int32)
	CtxGo(ctx context.Context, f func())
	SetPanicHandler(f func(context.Context, any))
	WorkCount() int32
}

var taskPool sync.Pool

func init() {
	taskPool.New = newTask
}

type task struct {
	ctx context.Context
	f   func()

	next *task
}

func (t *task) zero() {
	t.ctx = nil
	t.f = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.zero()
	taskPool.Put(t)
}

func newTask() any {
	return &task{}
}

type pool struct {
	name string

	cap      int32
	config   *Config
	taskHead *task
	taskTail *task
	taskLock sync.Mutex
	taskCnt  int32

	workerCnt int32

	panicHandler func(context.Context, any)
}

var _ Pool = (*pool)(nil)

func NewPool(name string, cap int32, config *Config) *pool {
	return &pool{
		name:   name,
		cap:    cap,
		config: config,
	}
}

func (p *pool) Name() string {
	return p.name
}
func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}
func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f
	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()
	atomic.AddInt32(&p.taskCnt, 1)
	if (atomic.LoadInt32(&p.taskCnt) >= p.config.ScaleThreshold && p.WorkCount() < atomic.LoadInt32(&p.cap)) || p.WorkCount() == 0 {
		p.incWorkerCount()
		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

func (p *pool) SetPanicHandler(f func(context.Context, any)) {
	p.panicHandler = f
}

func (p *pool) WorkCount() int32 {
	return atomic.LoadInt32(&p.workerCnt)
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCnt, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCnt, -1)
}
