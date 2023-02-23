package code

// 校验错误
const (
	ErrValidation int = iota + 10000
	ErrTokenValidation
)

// 数据库错误
const (
	ErrRecordNotFound int = iota + 10100
	ErrNameDuplicate
	ErrDatabaseOp
	ErrRecordNotEnough
	ErrNoUpdateRows
	ErrSystemErr
)

// 业务错误
const (
	ErrPassword int = iota + 10200
	ErrTokenExpired
	ErrAccessDenied
	ErrPlzLogin
)
