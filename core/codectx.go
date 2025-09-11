package core

import (
	"context"

	myErrors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/errors/code"
)

type (
	ctxCodeKey struct{}
	ctxMsgKey  struct{}
	ctxErrKey  struct{}
)

func CodeFrom(ctx context.Context) code.Code {
	if code, ok := ctx.Value(ctxCodeKey{}).(code.Code); ok {
		return code
	}
	return 0
}

func WithCode(ctx context.Context, code code.Code) context.Context {
	return context.WithValue(ctx, ctxCodeKey{}, code)
}

func MsgFrom(ctx context.Context) string {
	if code, ok := ctx.Value(ctxMsgKey{}).(string); ok {
		return code
	}
	return ""
}

func WithMsg(ctx context.Context, msg string) context.Context {
	return context.WithValue(ctx, ctxMsgKey{}, msg)
}

func ErrFrom(ctx context.Context) *myErrors.APIError {
	if err, ok := ctx.Value(ctxErrKey{}).(*myErrors.APIError); ok {
		return err
	}
	return nil
}

func WithErr(ctx context.Context, err *myErrors.APIError) context.Context {
	return context.WithValue(ctx, ctxErrKey{}, err)
}
