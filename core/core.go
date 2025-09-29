package core

import (
	"net/http"

	myErrors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/errors/code"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ListResponseData[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

type ApiResponse struct {
	ApiError *myErrors.APIError
	Data     interface{}
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	var (
		httpCode int = http.StatusOK
		resp         = Response{Data: data}
	)
	apiErr := ErrFrom(c)
	if apiErr != nil {
		httpCode = apiErr.StatusCode
		resp.Code = apiErr.Code
		resp.Message = apiErr.Message
	} else if err != nil {
		co := CodeFrom(c)
		if co == 0 {
			co = code.CodeSystemErr
		}
		resp.Code = int(co)
		msg := MsgFrom(c)
		if msg == "" {
			msg = co.Message()
		}
		resp.Message = msg
	}
	c.JSON(httpCode, resp)
}

func WriteListResponse(c *gin.Context, apiError *myErrors.APIError, total int64, data interface{}) {
	if apiError == nil {
		c.JSON(http.StatusOK, Response{Data: map[string]any{
			"total": total,
			"list":  data,
		}})
		return
	}

	c.JSON(apiError.StatusCode, Response{
		Code:    apiError.Code,
		Message: apiError.Message,
		Data:    data,
	})
}
