package core

import (
	"common/errors"
	myErrors "common/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApiResponse struct {
	ApiError *myErrors.APIError
	Data     interface{}
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *gin.Context, apiError *errors.APIError, data interface{}) {
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

func WriteListResponse(c *gin.Context, apiError *errors.APIError, total int64, data interface{}) {
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
