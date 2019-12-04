package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	INSTALLERROR        = &Errno{Code: -1, Message: "安装失败"}
	INSTALLED        = &Errno{Code: -1, Message: "你已经安装过数据库"}
	InternalServerError = &Errno{Code: 10001, Message: "接口服务器错误"}
	ErrBind             = &Errno{Code: 10002, Message: "绑定数据错误"}

	ErrValidation = &Errno{Code: 20001, Message: "验证失败"}
	ErrDatabase   = &Errno{Code: 20002, Message: "数据库出错"}

	// user errors
	ErrUserNotFound      = &Errno{Code: 20101, Message: "该用户不存在"}
	ErrPasswordIncorrect = &Errno{Code: 20102, Message: "密码错误"}
)
