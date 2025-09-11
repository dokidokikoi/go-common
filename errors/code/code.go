//go:generate stringer -type=Code -linecomment
package code

// Code 定义业务错误码
type Code int

// 客户端错误
const (
	CodeValidation      Code = iota + 10000 // 参数校验错误
	CodeTokenValidation                     // Token 校验失败
	CodePassword                            // 密码错误
	CodeTokenExpired                        // Token 过期
	CodeAccessDenied                        // 没有权限
	CodePlzLogin                            // 请先登录
)

// 数据库错误
const (
	CodeRecordNotFound  Code = iota + 10100 // 记录未找到
	CodeNameDuplicate                       // 名称重复
	CodeDatabaseOp                          // 数据库操作错误
	CodeRecordNotEnough                     // 记录数量不足
	CodeNoUpdateRows                        // 没有更新任何行
)

// 系统错误
const (
	CodeSystemErr Code = iota + 10200 // 系统内部错误
)

// codeMessages 用于返回业务友好的错误提示
var codeMessages = map[Code]string{
	CodeValidation:      "参数校验错误",
	CodeTokenValidation: "Token 校验失败",
	CodePassword:        "密码错误",
	CodeTokenExpired:    "Token 已过期",
	CodeAccessDenied:    "没有权限",
	CodePlzLogin:        "请先登录",

	CodeRecordNotFound:  "记录未找到",
	CodeNameDuplicate:   "名称重复",
	CodeDatabaseOp:      "数据库操作错误",
	CodeRecordNotEnough: "记录数量不足",
	CodeNoUpdateRows:    "没有更新任何行",

	CodeSystemErr: "系统内部错误",
}

// Message 返回错误码对应的提示信息
func (c Code) Message() string {
	if msg, ok := codeMessages[c]; ok {
		return msg
	}
	return "未知错误"
}
