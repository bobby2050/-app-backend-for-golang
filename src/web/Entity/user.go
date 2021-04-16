package Entity

// 绑定为json
type ReqUserLogin struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password"  binding:"required"`
}

type ReqUserRegistion struct {
	UserName   string `form:"username" json:"username" binding:"required"`
	Password   string `form:"password" json:"password"  binding:"required"`
	RePassword string `form:"re_password" json:"re_password"  binding:"required"`
	NickName   string `form:"nick_name" json:"nick_name" binding:"required"`
	Email      string `form:"email" json:"email" binding:"required"`
}

// 接收客户端的文章内容
type ReqPublish struct {
	Token          string `form:"token" json:"token" binding:"required"`
	CompanyName    string `form:"companyName" json:"companyName"  binding:"required"`
	DepartmentName string `form:"departmentName" json:"departmentName"  binding:"required"`
	GroupName      string `form:"groupName" json:"groupName"  binding:"required"`
	PositionName   string `form:"positionName" json:"positionName"  binding:"required"`
	ProvinceName   string `form:"provinceName" json:"provinceName"  binding:"required"`
	CityName       string `form:"cityName" json:"cityName"  binding:"required"`
	ContentText    string `form:"contentText" json:"contentText"  binding:"required"`
}

// 接收客户端的用户喜欢文章
type ReqUserLikePublish struct {
	Token  string `form:"token" json:"token" binding:"required"`
	PostId int64 `form:"post_id" json:"post_id"  binding:"required"`
	State  int    `form:"state" json:"state" `
}

// 接收客户端的用户关注列表
type ReqUserFocusPublishList struct {
	Flag      string `form:"flag" json:"flag"  binding:"required"`
	Token     string `form:"token" json:"token"`
	PublishId int    `form:"publishId" json:"publishId"`
}

// 文章回复列表
type ReqPublishReplyList struct {
	PublishId int    `form:"publish_id" json:"publish_id"`
}

// 接收客户端的用户
type ReqToken struct {
	Token     string `form:"token" json:"token"`

}

// 关注Publish
type ReqFocusPublish struct {
	Token  string `form:"token" json:"token"  binding:"required"`
	PublishId int    `form:"publish_id" json:"publish_id"  binding:"required"`
}

// 回复内容列表
type ReqPublishReply struct {
	Token     string `form:"token" json:"token" `
	PublishId int    `form:"publishId" json:"publishId"  binding:"required"`
	Page      int    `form:"page" json:"page"  binding:"required"`
}

// 回复内容
type Reply struct {
	PublishId int    `form:"publish_id" json:"publish_id"  binding:"required"`
	Token     string `form:"token" json:"token" `
	ReplyText string `form:"reply_text" json:"reply_text"  binding:"required"`
}
// 我的关注
type MyFocusUser struct {
	Token  string `form:"token" json:"token"  binding:"required"`
}

// 我的发布列表
type MyPublish struct {
	Token  string `form:"token" json:"token"  binding:"required"`
}



// 更新coredata的回复数和喜欢数
type PostReplyLike struct {
	Token string `form:"token" json:"token"`
	PostId []  int64  `form:"post_ids" json:"post_ids" binding:"required"`

}

// 举报
type ReportPost struct {
	Token  string `form:"token" json:"token"`
	ReportType  int `form:"reportType" json:"reportType"  binding:"required"`
	PostId  int64 `form:"postId" json:"postId"  binding:"required"`
}