package middleware

import (
	"context"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PreHandle[Input any, Resp any](handler func(ctx context.Context, input *Input) (Resp, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input Input
		err := ctx.ShouldBind(&input)
		if err != nil {
			zaplog.L().Error("ShouldBind", zap.Error(err))
			core.WithErr(ctx, errors.ApiErrValidation)
			core.WriteResponse(ctx, err, nil)
			return
		}

		resp, err := handler(ctx, &input)
		if err != nil {
			zaplog.L().Error("handler", zap.Error(err))
			core.WriteResponse(ctx, err, nil)
			return
		}
		core.WriteResponse(ctx, nil, resp)
	}
}
