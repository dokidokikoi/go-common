package gopool

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

var workerPool sync.Pool

func init() {
	workerPool.New = newWorker
}

type worker struct {
	pool *pool
}

func newWorker() any {
	return &worker{}
}

func (w *worker) run() {
	go func() {
		for {
			var t *task
			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCnt, -1)
			}
			if t == nil {
				w.close()
				w.pool.taskLock.Unlock()
				w.Recycle()
				return
			}
			w.pool.taskLock.Unlock()

			func() {
				defer func() {
					if r := recover(); r != nil {
						if w.pool.panicHandler != nil {
							w.pool.panicHandler(t.ctx, r)
						} else {
							msg := fmt.Sprintf("GOPOOL: panic in pool: %s: %v: %s", w.pool.name, r, debug.Stack())
							zaplog.L().Error(msg)
						}
					}
				}()
				t.f()
			}()
			t.Recycle()
		}
	}()
}

func (w *worker) close() {
	w.pool.decWorkerCount()
}

func (w *worker) zero() {
	w.pool = nil
}

func (w *worker) Recycle() {
	w.zero()
	workerPool.Put(w)
}
