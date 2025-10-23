package middleware

import (
	"context"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/errors/code"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PreHandleOptions struct {
	co  code.Code
	msg string
	err *errors.APIError
}

func (m *PreHandleOptions) SetCode(co code.Code) {
	m.co = co
}

func (m *PreHandleOptions) SetMsg(msg string) {
	m.msg = msg
}

func (m *PreHandleOptions) SetErr(err *errors.APIError) {
	m.err = err
}

func (m *PreHandleOptions) ToResponseOpt() core.ResponseOption {
	return core.ResponseOption{
		Code: m.co,
		Err:  m.err,
		Msg:  m.msg,
	}
}

func PreHandle[Input any, Resp any](handler func(ctx context.Context, input *Input, op *PreHandleOptions) (Resp, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := zaplog.From(ctx)
		var input Input
		err := ctx.ShouldBind(&input)
		if err != nil {
			logger.Error("ShouldBind", zap.Error(err))
			core.WriteResponse(ctx, err, nil)
			return
		}

		op := &PreHandleOptions{}
		resp, err := handler(ctx, &input, op)
		if err != nil {
			logger.Error("handler", zap.Error(err))
			core.WriteResponse(ctx, err, nil, op.ToResponseOpt())
			return
		}
		core.WriteResponse(ctx, nil, resp)
	}
}
