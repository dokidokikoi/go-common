package middleware

import (
	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func PreHandle[Input any, Resp any](handler func(ctx *gin.Context, input Input) (Resp, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input Input
		err := ctx.ShouldBind(&input)
		if err != nil {
			core.WriteResponse(ctx, errors.ApiErrValidation, nil)
			return
		}

		resp, err := handler(ctx, input)
		if err != nil {
			core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
			return
		}
		core.WriteResponse(ctx, nil, resp)
	}
}
