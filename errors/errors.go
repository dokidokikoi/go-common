package errors

import (
	"errors"

	"github.com/dokidokikoi/go-common/errors/code"
)

var (
	ErrRecordNotFound  = errors.New("该记录未找到")
	ErrValidation      = errors.New("参数格式错误")
	ErrDatabaseOp      = errors.New("录入失败")
	ErrRecordNotEnough = errors.New("记录数不足")
	ErrNameDuplicate   = errors.New("名称重复")
	ErrNoUpdateRows    = errors.New("无更新记录")
	ErrSystemErr       = errors.New("系统错误")
	ErrTokenValidation = errors.New("token有误")
	ErrAccessDenied    = errors.New("拒绝访问")
	ErrTokenExpired    = errors.New("token已过期")
	ErrPassword        = errors.New("密码错误")
	ErrPlzLogin        = errors.New("请登录")
)

var (
	ApiErrRecordNotFound  = NotifyFailed(ErrRecordNotFound.Error(), code.ErrRecordNotFound)
	ApiErrValidation      = ClientFailed(ErrValidation.Error(), code.ErrValidation)
	ApiErrDatabaseOp      = ClientFailed(ErrDatabaseOp.Error(), code.ErrDatabaseOp)
	ApiErrRecordNotEnough = ClientFailed(ErrRecordNotEnough.Error(), code.ErrRecordNotEnough)
	ApiErrNameDuplicate   = ClientFailed(ErrNameDuplicate.Error(), code.ErrNameDuplicate)
	ApiErrNoUpdateRows    = ClientFailed(ErrNoUpdateRows.Error(), code.ErrNoUpdateRows)
	ApiErrSystemErr       = ClientFailed(ErrSystemErr.Error(), code.ErrSystemErr)
	ApiErrTokenValidation = ClientFailed(ErrTokenValidation.Error(), code.ErrTokenValidation)
	ApiErrAccessDenied    = ClientFailed(ErrAccessDenied.Error(), code.ErrAccessDenied)
	ApiErrTokenExpired    = ClientFailed(ErrTokenExpired.Error(), code.ErrTokenExpired)
	ApiErrPassword        = ClientFailed(ErrPassword.Error(), code.ErrPassword)
	ApiErrPlzLogin        = ClientFailed(ErrPlzLogin.Error(), code.ErrPlzLogin)
)
