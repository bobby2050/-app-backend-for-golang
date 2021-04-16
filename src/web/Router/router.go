package Router

import (
	"github.com/gin-gonic/gin"
	"meizi/src/web/Cache"
	"meizi/src/web/Middleware"
	"meizi/src/web/Model"
	"meizi/src/web/Service"
)

func Route(r *gin.Engine)  {

	// 设置静态资源目录
	r.Static("/assets", "./assets")

	// 检测方法
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// v1 版本
	v1 := r.Group("/v1")

	// 登陆
	v1.POST("/login",  Service.Login)

	// 注册
	v1.POST("/register",  Service.Register)

	// 请求推荐列表
	//v1.GET("/list",  Service.PublishList)

	// 请求关注列表
	//v1.GET("/focusList",  Service.FocusPublishList)

	// 请求喜欢文章
	//v1.POST("/like",  Service.LikePublish)

	// 请求Publish
	//v1.POST("/focusPublish",  Service.FocusPublish)

	// 请求取消关注用户
	//v1.POST("/cancelFocusUser",  Service.CancelFocusUser)

	// 请求回复列表
	//v1.POST("/publishReply",  Service.PublishReply)

	// 回复
	v1.POST("/reply",  Service.Reply)

	// 新增文章
	v1.POST("/publishPost",  Service.Write)

	// 我的关注
	//v1.POST("/myfocus",  Service.MyFocus)

	// 我的发布列表
	//v1.POST("/mypublish",  Service.MyPublish)

	// 接收服务器最新的通知
	v1.POST("/serverNotication",  Service.ServerNotication)


	// --------------------------- v1 版本 redis版 ---------------------------
	v2 := r.Group("/cache")

	// 请求推荐列表
	v2.GET("/InitPostlist",  Model.InitPostToRedis)
	v2.GET("/InitFocuslist",  Model.InitFocusToRedis)
	v2.GET("/InitUserLovePublishlist",  Model.InitUserLovePublishToRedis)
	v2.GET("/InitPublishReplylist",  Model.InitPublishReplyToRedis)

	// 请求推荐列表
	v2.GET("/list",  Cache.PublishList)

	// 当前用户喜欢的文章
	v2.GET("/myPublishLike",  Cache.MyPublishLike)

	// 获取用户关注的文章编号
	v2.POST("/myfocusPostId",  Cache.MyfocusPostId)

	// 请求关注列表
	v2.GET("/listFocusPublish",  Cache.FocusPublishList)

	// 请求回复列表
	v2.POST("/publishReply",  Cache.PublishReply)

	// 请求喜欢文章
	v2.POST("/like",  Cache.LikePublish)

	// 回复
	v2.POST("/reply",  Cache.Reply)

	// 更新coredata的回复数和喜欢数
	v2.POST("/postReplyLikeCount",  Cache.PostReplyLikeCount)

	// 关注文章
	v2.POST("/focusPublish",  Cache.FocusPublish)

	// 取消关注文章
	v2.POST("/cancelFocusPublish",  Cache.CancelFocusPublish)

	// 我的发布列表
	v2.POST("/mypublish",  Cache.MyPublish)

	// 举报
	v2.POST("/reprotPost",  Cache.ReprotPost)

	//// 请求喜欢文章
	//v2.POST("/like",  Service.LikePublish)
	//
	//// 请求关注用户
	//v2.POST("/focusUser",  Service.FocusUser)
	//
	//// 请求取消关注用户
	//v2.POST("/cancelFocusUser",  Service.CancelFocusUser)
	//
	//// 请求回复列表
	//v2.POST("/publishReply",  Service.PublishReply)
	//
	//// 回复
	//v2.POST("/reply",  Service.Reply)
	//
	//// 回复
	//v2.POST("/publishPost",  Service.Write)
	//
	//// 我的关注
	//v2.POST("/myfocus",  Service.MyFocus)
	//
	//// 我的发布列表
	//v2.POST("/mypublish",  Service.MyPublish)

	// 用户校验接口
	//r.Use(Middleware.CheckToken(), Middleware.Cors())
	v1.Use(Middleware.CheckToken())
	{

		v1.POST("/UserDetail", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "read",
			})
		})



	}





}
