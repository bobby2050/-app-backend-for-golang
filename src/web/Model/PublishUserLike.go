package Model

import "time"

//  发布内容
type PublishUserLike struct {
	Id        int       `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"` // 主键编号
	PublishId int       `json:"publish_id" gorm:"Column:publish_id"`    // 发布编号
	UserId    int       `json:"user_id" gorm:"Column:user_id"`          // 用户编号
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`           // 创建时间
}
