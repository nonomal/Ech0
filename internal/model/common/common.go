package model

// UserStatus 用于存储用户状态信息
type UserStatus struct {
	UserID   uint   `json:"user_id"`  // 用户ID
	UserName string `json:"username"` // 用户名
	IsAdmin  bool   `json:"is_admin"` // 是否是管理员
}

// Status 用于存储Echo状态信息
type Status struct {
	SysAdminID uint         `json:"sys_admin_id"` // 系统管理员ID
	Username   string       `json:"username"`     // 系统管理员用户名
	Users      []UserStatus `json:"users"`        // 所有用户
	Logo       string       `json:"logo"`         // 站点logo
	TotalEchos int          `json:"total_echos"`  // 总共发布数量
}

// TempFile 使用延迟删除机制处理S3和本地存储的孤儿文件
type TempFile struct {
	ID             uint   `json:"id" gorm:"primaryKey"` // 主键ID
	FileName       string `json:"file_name"`            // 文件名
	Storage        string `json:"storage"`              // 存储类型 local/s3/r2
	FileType       string `json:"file_type"`            // 文件类型 image/audio
	Bucket         string `json:"bucket"`               // 存储桶
	ObjectKey      string `json:"object_key"`           // 对象键
	Deleted        bool   `json:"deleted"`              // 是否已删除
	CreatedAt      int64  `json:"created_at"`           // 创建时间（Unix时间戳）
	LastAccessedAt int64  `json:"last_accessed_at"`     // 最后访问时间（Unix时间戳）
}

// Heatmap 用于存储热力图数据
type Heatmap struct {
	Date  string `json:"date"`  // 日期
	Count int    `json:"count"` // Echo数量
}

// File 相关

type UploadFileType string
type FileStorageType string
type CommentProvider string
type S3Provider string

const (
	// ImageType  图片类型
	ImageType UploadFileType = "image"
	// AudioType  音频类型
	AudioType UploadFileType = "audio"
)

const (
	// LOCAL_FILE 本地存储类型
	LOCAL_FILE FileStorageType = "local"
	// S3_FILE   S3 存储类型
	S3_FILE FileStorageType = "s3"
	// R2_FILE   R2 存储类型
	R2_FILE FileStorageType = "r2"
)

const (
	// Twikoo 评论服务
	TWIKOO CommentProvider = "twikoo"
	// Artalk 评论服务
	ARTALK CommentProvider = "artalk"
	// WALINE 评论服务
	WALINE CommentProvider = "waline"
	// GISCUS 评论服务
	GISCUS CommentProvider = "giscus"
)

const (
	AWS     S3Provider = "aws"
	ALIYUN  S3Provider = "aliyun"
	TENCENT S3Provider = "tencent"
	MINIO   S3Provider = "minio"
	OTHER   S3Provider = "other"
)

// key value表
type KeyValue struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Value string `json:"value"`
}

// 键值对相关
const (
	// SystemSettingsKey 是系统设置的键
	SystemSettingsKey = "system_settings"
	// CommentSettingKey 是评论设置的建
	CommentSettingKey = "comment_setting"
	// S3SettingKey 是 S3 存储设置的键
	S3SettingKey = "s3_setting"
	// ServerURLKey 是服务器URL设置的键
	ServerURLKey = "server_url"
	// MigrationKey 是数据库迁移的标记键
	MigrationKey = "db_migration:message_to_echo:v1"
)

// PageQueryResult 用于分页查询的结果数据传输对象
type PageQueryResult[T any] struct {
	Total int64 `json:"total"`
	Items T     `json:"items"`
}

const (
	// InitInstallCode 是初始化安装的标志
	InitInstallCode = 666
)

const (
	// Version 是当前版本号
	Version = "2.6.8"
)
