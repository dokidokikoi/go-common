package middleware

import (
	"context"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PreHandle[Input any, Resp any](handler func(ctx context.Context, input *Input) (Resp, *errors.APIError)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input Input
		err := ctx.ShouldBind(&input)
		if err != nil {
			zaplog.L().Error("ShouldBind", zap.Error(err))
			core.WriteResponse(ctx, errors.ApiErrValidation, nil)
			return
		}

		resp, apiErr := handler(ctx, &input)
		if apiErr != nil {
			zaplog.L().Error("handler", zap.Error(err))
			core.WriteResponse(ctx, apiErr, nil)
			return
		}
		core.WriteResponse(ctx, nil, resp)
	}
}
