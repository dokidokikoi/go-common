package code

// 客户端错误
const (
	ErrValidation int = iota + 10000
	ErrTokenValidation
	ErrPassword
	ErrTokenExpired
	ErrAccessDenied
	ErrPlzLogin
)

// 数据库错误
const (
	ErrRecordNotFound int = iota + 10100
	ErrNameDuplicate
	ErrDatabaseOp
	ErrRecordNotEnough
	ErrNoUpdateRows
)

// 系统错误
const (
	ErrSystemErr int = iota + 10200
)
