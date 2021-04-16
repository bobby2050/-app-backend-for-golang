package Entity

//  发布内容
type RetPublish struct {
	Id          int    `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"`  // 主键编号
	UserId      int    `json:"user_id" gorm:"Column:user_id"`           // 用户编号
	UserName    string `json:"user_name" gorm:"Column:user_name"`       // 用户名
	NickName    string `json:"nick_name" gorm:"Column:nick_name"`       // 昵称
	PublishText string `json:"publish_text" gorm:"Column:publish_text"` // 发布内容
	Status      int    `json:"status" gorm:"Column:status"`             // 审核状态 0 未审核  1 审核
	ReplyCount  int64  `json:"reply_count" gorm:"Column:reply_count"`   // 回复数
	LikeCount   int64  `json:"like_count" gorm:"Column:like_count"`     // 喜欢数
	CompanyName string `json:"company_name" gorm:"Column:company_name"` // 公司名称
	CreatedAt   string `json:"created_at" gorm:"created_at"`            // 创建时间
	UpdatedAt   string `json:"updated_at" gorm:"updated_at"`            // 更新时间
	//IsLike      bool      `json:"is_like" gorm:"Column:is_like"`           // 当前用户喜欢
	//IsFocus     bool      `json:"is_focus" gorm:"Column:is_focus"`         // 是否关注用户

	DepartmentName string `json:"department_name" gorm:"Column:department_name"` // 部门
	GroupName      string `json:"group_name" gorm:"Column:group_name"`           // 组名
	PositionName   string `json:"position_name" gorm:"Column:position_name"`     // 职位
	ProvinceName   string `json:"province_name" gorm:"Column:province_name"`     // 省份
	CityName       string `json:"city_name" gorm:"Column:city_name"`             // 城市
}

type RetPublishReply struct {
	Id        int    `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"` // 主键编号
	PublishId int    `json:"publish_id" gorm:"Column:publish_id"`    // 内容编号
	UserId    int    `json:"user_id" gorm:"Column:user_id"`          // 用户编号
	NickName  string `json:"nick_name" gorm:"Column:nick_name"`      // 昵称
	ReplyText string `json:"reply_text" gorm:"Column:reply_text"`    // 内容
	Status    int    `json:"status" gorm:"Column:status"`            // 审核状态 0 未审核  1 审核
	CreatedAt string `json:"created_at" gorm:"created_at"`           // 创建时间
}
