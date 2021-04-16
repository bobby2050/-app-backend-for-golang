package Model

import "time"

type User struct {
	Id             int       `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"` // 主键编号
	UserName       string    `json:"user_name" gorm:"Column:username"`       // 用户名
	Password       string    `json:"password" gorm:"Column:password"`        // 密码
	Salt           string    `json:"salt" gorm:"Column:salt"`                // 盐值
	NickName       string    `json:"nick_name" gorm:"Column:nick_name"`      // 昵称
	Mobile         string    `json:"mobile" gorm:"Column:mobile"`            // 手机号
	Email          string    `json:"email" gorm:"Column:email"`              // 邮件
	LastSessionKey string    `json:"session" gorm:"Column:last_session_key"` // 最后登陆的sessionid
	CreatedAt      time.Time `json:"created_at" gorm:"created_at"`           // 创建时间
	UpdatedAt      time.Time `json:"updated_at" gorm:"updated_at"`           // 更新时间
}
