package core

import (
	"net/http"

	"github.com/dokidokikoi/go-common/errors"
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

type ResponseOption struct {
	Code code.Code
	Msg  string
	Err  *myErrors.APIError
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *gin.Context, err error, data interface{}, opt ...ResponseOption) {
	var (
		httpCode int = http.StatusOK
		resp         = Response{Data: data}
	)
	resp.Data = data

	var (
		apiErr *errors.APIError
		co     code.Code
		msg    string
	)
	if len(opt) > 0 {
		apiErr = opt[0].Err
		co = opt[0].Code
		msg = opt[0].Msg
	}
	if apiErr != nil {
		httpCode = apiErr.StatusCode
		resp.Code = apiErr.Code
		resp.Message = apiErr.Message
	} else if err != nil {
		if co == 0 {
			co = code.CodeSystemErr
		}
		resp.Code = int(co)
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
