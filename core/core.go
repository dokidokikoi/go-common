package core

import (
	"net/http"

	myErrors "github.com/dokidokikoi/go-common/errors"

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
func WriteResponse(c *gin.Context, apiError *myErrors.APIError, data interface{}) {
	if apiError == nil {
		c.JSON(http.StatusOK, Response{Data: data})
		return
	}

	c.JSON(apiError.StatusCode, Response{
		Code:    apiError.Code,
		Message: apiError.Message,
		Data:    data,
	})
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
