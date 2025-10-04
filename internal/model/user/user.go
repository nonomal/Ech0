package model

const (
	USER_NOT_EXISTS_ID = 0
)

// User 定义用户实体
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	IsAdmin  bool   `gorm:"bool" json:"is_admin"`
	Avatar   string `gorm:"size:255" json:"avatar"`

	// ActivityPub 相关字段
	// DisplayName   string    `gorm:"size:128" json:"display_name"`
	// ActorURL      string    `gorm:"size:256;unique" json:"actor_url"` // Actor 对象地址
	// PublicKeyPEM  string    `gorm:"type:text" json:"public_key_pem"`  // 公钥
	// PrivateKeyPEM string    `gorm:"type:text" json:"-"`               // 私钥
	// CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	// UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
}
