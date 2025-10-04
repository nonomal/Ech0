package model

// ServerError 定义服务器错误信息
type ServerError struct {
	Msg string
	Err error
}

// 失败相关的常量
const (
	INVALID_FILE_PATH      = "无效的文件路径"
	INVALID_REQUEST_BODY   = "无效的请求体"
	INVALID_PARAMS_BODY    = "无效参数"
	INVALID_QUERY_PARAMS   = "无效的查询参数"
	INVALID_REQUEST_METHOD = "无效的请求方法"
)

// Auth 错误相关常量
const (
	USERNAME_OR_PASSWORD_NOT_BE_EMPTY = "用户名或密码不能为空"
	PASSWORD_INCORRECT                = "密码错误"
	USER_NOTFOUND                     = "用户不存在"
	USER_COUNT_EXCEED_LIMIT           = "用户数量超过限制"
	USERNAME_HAS_EXISTS               = "用户名已存在"
	TOKEN_NOT_FOUND                   = "未找到令牌,请点击右上角登录"
	TOKEN_NOT_VALID                   = "令牌无效，请重新登录"
	TOKEN_PARSE_ERROR                 = "令牌解析失败，请尝试重新登陆"
	USER_REGISTER_NOT_ALLOW           = "当前系统禁止注册新用户"
)

// Echo 错误相关常量
const (
	NO_PERMISSION_DENIED  = "没有权限,请联系系统管理员"
	ECHO_CAN_NOT_BE_EMPTY = "ECHO 内容不能为空"
	ECHO_NOT_FOUND        = "找不到Echo"
)

// Common 错误相关常量
const (
	NO_FILE_UPLOAD_ERROR   = "找不到上传的文件"
	NO_FILE_STORAGE_ERROR  = "未知存储方式"
	FILE_TYPE_NOT_ALLOWED  = "不支持的文件类型"
	FILE_SIZE_EXCEED_LIMIT = "文件大小超过限制"
	IMAGE_NOT_FOUND        = "图片未找到"
	INVALID_PARAMS         = "错误的参数"
	SIGNUP_FIRST           = "请先注册用户"
	S3_NOT_ENABLED         = "S3存储未启用"
	S3_NOT_CONFIGURED      = "S3存储未配置"
	S3_CONFIG_ERROR        = "S3存储配置错误"
)

// User 错误相关常量
const (
	USERNAME_ALREADY_EXISTS = "用户名已存在"
)

// TO DO 错误相关常量
const (
	TODO_EXCEED_LIMIT = "待办事项数量已达上限"
)

// Connect 错误相关常量
const (
	INVALID_CONNECTION_URL = "connect url不能为空"
	CONNECT_HAS_EXISTS     = "connect 已经存在"
)

// Setting 错误相关常量
const (
	NO_SUCH_COMMENT_PROVIDER = "无效的评论服务提供者"
)

// Backup 错误相关常量
const (
	SNAPSHOT_UPLOAD_FAILED  = "快照上传失败"
	SNAPSHOT_RESTORE_FAILED = "快照恢复失败"
	DATABASE_CLOSE_FAILED   = "数据库关闭失败"
)

// Fediverse 错误相关常量
const (
	GET_ACTOR_ERROR         = "获取 Actor 信息失败"
	ACTIVEPUB_NOT_ENABLED   = "ActivityPub 未启用"
	FEDIVERSE_INVALID_INPUT = "无效的联邦参数"
	FOLLOW_RELATION_MISSING = "未找到关注关系"
)
