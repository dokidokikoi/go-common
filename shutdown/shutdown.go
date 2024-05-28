package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var sign struct {
	c   chan os.Signal
	ctx context.Context
}

func init() {
	sign = struct {
		c   chan os.Signal
		ctx context.Context
	}{
		c:   make(chan os.Signal, 1),
		ctx: context.Background(),
	}
}

func WithSignals(signals ...syscall.Signal) {
	for _, s := range signals {
		signal.Notify(sign.c, s)
	}
}

func WithContext(ctx context.Context) {
	sign.ctx = ctx
}

func Close(funcs ...func()) {
	select {
	case <-sign.c:
	case <-sign.ctx.Done():
	}
	signal.Stop(sign.c)

	for _, f := range funcs {
		f()
	}
}
