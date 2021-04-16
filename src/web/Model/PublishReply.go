package Model

import "time"

//  发布内容
type PublishReply struct {
	Id        int       `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"` // 主键编号
	PublishId int       `json:"publish_id" gorm:"Column:publish_id"`    // 内容编号
	UserId    int       `json:"user_id" gorm:"Column:user_id"`          // 用户编号
	NickName  string    `json:"nick_name" gorm:"Column:nick_name"`      // 昵称
	ReplyText string    `json:"reply_text" gorm:"Column:reply_text"`    // 内容
	Status    int       `json:"status" gorm:"Column:status"`            // 审核状态 0 未审核  1 审核
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`           // 创建时间
}
