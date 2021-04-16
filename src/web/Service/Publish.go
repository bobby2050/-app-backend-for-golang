package Service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"meizi/src/Db"
	"meizi/src/Redis"
	"meizi/src/web/Entity"
	"meizi/src/web/Libaries"
	"meizi/src/web/Model"
	"strconv"
	"time"
)

var ctx = context.Background()

func ServerNotication(c *gin.Context) {
	res, err := Redis.RedisPool.Get(ctx,"ServerNotication").Result()
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": res,
	})
}

// 发布内容
func Write(c *gin.Context) {
	fmt.Println("发布内容")
	var reqPublish Entity.ReqPublish
	if c.BindJSON(&reqPublish) == nil { // 参数校验正确
		fmt.Println(reqPublish)

		//val, err := Redis.RedisPool.HMGet(ctx, reqPublish.Token,  "user_id", "nick_name", "username").Result()
		//if err != nil {
		//	panic(err)
		//}
		redisUserInfo := Libaries.GetUserIdByToken(reqPublish.Token)
		fmt.Println("key", redisUserInfo)
		dbUserId := redisUserInfo[0]
		nickName := redisUserInfo[1]
		username := redisUserInfo[2]

		userId, _ := strconv.Atoi(fmt.Sprintf("%v", dbUserId))
		myPublishText := &Model.Publish{
			UserName:       fmt.Sprintf("%v", username),
			UserId:         userId,
			NickName:       fmt.Sprintf("%v", nickName),
			PublishText:    reqPublish.ContentText,
			Status:         0,
			CreatedAt:      time.Now(),
			CompanyName:    reqPublish.CompanyName,
			DepartmentName: reqPublish.DepartmentName,
			GroupName:      reqPublish.GroupName,    // 组名
			PositionName:   reqPublish.PositionName, // 职位
			ProvinceName:   reqPublish.ProvinceName, // 省份
			CityName:       reqPublish.CityName,     // 城市
		}
		Db.DBHelper.Create(myPublishText)

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": "",
		})
	}

}

// 关注列表
func FocusPublishList(c *gin.Context) {
	fmt.Println("-----------------------------------------")
	fmt.Println("关注列表")
	fmt.Println("-----------------------------------------")

	var reqUserFocusPublishList Entity.ReqUserFocusPublishList
	if c.Bind(&reqUserFocusPublishList) == nil { // 参数校验正确
		flag := reqUserFocusPublishList.Flag // 获取Get参数 page
		fmt.Println("flag: " + flag)

		//token := c.Request.Header.Get("token") // 获取Get参数 token
		//fmt.Println("token: " + token)
		token := reqUserFocusPublishList.Token // 获取Get参数 token
		fmt.Println("token： " + token)

		redisUserInfo := Libaries.GetUserIdByToken(token)
		fmt.Println("key", redisUserInfo)
		db_user_id := redisUserInfo[0]
		fmt.Println("user_id：", db_user_id)
		//userId := val[0]
		//fmt.Println("userId:")
		//fmt.Println(userId)

		// 2.查找出关注的用户id
		var rserFocusPublish []Model.UserFocusPublish
		Db.DBHelper.Where("user_id", db_user_id).Find(&rserFocusPublish)

		var focus []int
		for i := 0; i < len(rserFocusPublish); i++ {
			focus = append(focus, rserFocusPublish[i].PublishId)
		}

		id := reqUserFocusPublishList.PublishId

		var appendsql = ""
		if flag == "first" {
			appendsql = " AND p.id > " + strconv.Itoa(id)
		} else if flag == "last" {
			appendsql = " AND p.id < " + strconv.Itoa(id)
		}

		// 3.查找出关注用户id的发布列表
		//Db.DBHelper.Where("status = ? AND user_id in(?)", 1, focus).Find(&publishs)
		//var publishs []Entity.Publish
		var sql = ""

		sql = "SELECT p.*, IF(l.id>0,1,0) as is_like, 0 as is_focus  " +
			"FROM tbl_publish p \nLEFT JOIN \ntbl_publish_user_like l ON (l.publish_id=p.id AND l.user_id = ?) " +
			"WHERE p.status=1 AND p.user_id in (?) " + appendsql + " ORDER BY p.id DESC LIMIT 5"

		var publishs []Entity.Publish
		Db.DBHelper.Raw(sql, db_user_id, focus).Scan(&publishs)
		// 查找出关注用户id的发布列表

		fmt.Println(publishs)

		timeLayout := "2006-01-02 15:04:05" //转化所需模板
		//loc, _ := time.LoadLocation("Local")    //获取时区

		var retPublish []Entity.RetPublish
		for i := 0; i < len(publishs); i++ {
			retPublish = append(retPublish, Entity.RetPublish{
				Id:          publishs[i].Id,
				UserId:      publishs[i].UserId,
				UserName:    publishs[i].UserName,
				NickName:    publishs[i].NickName,
				PublishText: publishs[i].PublishText,
				Status:      publishs[i].Status,
				//ReplyCount:     publishs[i].ReplyCount,
				//LikeCount:      publishs[i].LikeCount,
				CompanyName: publishs[i].CompanyName,
				CreatedAt:   publishs[i].CreatedAt.Format(timeLayout),
				UpdatedAt:   publishs[i].UpdatedAt.Format(timeLayout),
				//IsLike:         publishs[i].IsLike,
				//IsFocus:        publishs[i].IsFocus,
				DepartmentName: publishs[i].DepartmentName, //
				GroupName:      publishs[i].GroupName,      //string    `json:"group_name" gorm:"Column:group_name"`           // 组名
				PositionName:   publishs[i].PositionName,   //string    `json:"position_name" gorm:"Column:position_name"`     // 职位
				ProvinceName:   publishs[i].ProvinceName,   //string    `json:"province_name" gorm:"Column:province_name"`     // 省份
				CityName:       publishs[i].CityName,       //string    `json:"city_name" gorm:"Column:city_name"`             // 城市
			})
		}

		//fmt.Println(tc)
		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": retPublish,
		})

	}

	fmt.Println("-----------------------------------------")
	fmt.Println("关注列表结束")
	fmt.Println("-----------------------------------------")

}

// 推荐列表
func PublishList(c *gin.Context) {
	fmt.Println("-----------------------------------------")
	fmt.Println("推荐列表")
	fmt.Println("-----------------------------------------")

	var reqUserFocusPublishList Entity.ReqUserFocusPublishList
	if c.Bind(&reqUserFocusPublishList) == nil { // 参数校验正确
		flag := reqUserFocusPublishList.Flag // 获取Get参数 page
		fmt.Println("flag: " + flag)

		id := reqUserFocusPublishList.PublishId

		var appendsql = ""
		if flag == "first" {
			appendsql = " AND p.id > " + strconv.Itoa(id)
		} else if flag == "last" {
			appendsql = " AND p.id < " + strconv.Itoa(id)
		}

		token := reqUserFocusPublishList.Token // 获取Get参数 token
		fmt.Println("token： " + token)

		redisUserInfo := Libaries.GetUserIdByToken(token)
		fmt.Println("key", redisUserInfo)
		db_user_id := redisUserInfo[0]
		fmt.Println("user_id：", db_user_id)
		//nick_name := redisUserInfo[1]
		//username := redisUserInfo[2]

		var publishs []Entity.Publish
		var sql = ""
		if db_user_id == "" {
			sql = "SELECT p.*, 0 as is_like, 0 as is_focus  FROM tbl_publish p  WHERE   p.status=1 " + appendsql + " ORDER BY p.id DESC Limit 10"

		} else {
			sql = "SELECT p.*, IF(l.id>0,1,0) as is_like, " +
				"IF(f.focus_user_id>0,1,0) AS is_focus " +
				"FROM tbl_publish p LEFT JOIN tbl_publish_user_like l ON (l.publish_id=p.id AND l.user_id = ?) " +
				"LEFT JOIN tbl_user_focus f ON (f.focus_user_id = p.user_id AND f.user_id = ?) " +
				"WHERE p.status=1 " + appendsql + " ORDER BY p.id DESC Limit 5"

		}
		Db.DBHelper.Raw(sql, db_user_id, db_user_id).Scan(&publishs)
		// 查找出关注用户id的发布列表

		timeLayout := "2006-01-02 15:04:05" //转化所需模板
		//loc, _ := time.LoadLocation("Local")    //获取时区

		var retPublish []Entity.RetPublish
		for i := 0; i < len(publishs); i++ {
			retPublish = append(retPublish, Entity.RetPublish{
				Id:          publishs[i].Id,
				UserId:      publishs[i].UserId,
				UserName:    publishs[i].UserName,
				NickName:    publishs[i].NickName,
				PublishText: publishs[i].PublishText,
				Status:      publishs[i].Status,
				//ReplyCount:     publishs[i].ReplyCount,
				//LikeCount:      publishs[i].LikeCount,
				CompanyName: publishs[i].CompanyName,
				CreatedAt:   publishs[i].CreatedAt.Format(timeLayout),
				UpdatedAt:   publishs[i].UpdatedAt.Format(timeLayout),
				//IsLike:         publishs[i].IsLike,
				//IsFocus:        publishs[i].IsFocus,
				DepartmentName: publishs[i].DepartmentName, //
				GroupName:      publishs[i].GroupName,      //string    `json:"group_name" gorm:"Column:group_name"`           // 组名
				PositionName:   publishs[i].PositionName,   //string    `json:"position_name" gorm:"Column:position_name"`     // 职位
				ProvinceName:   publishs[i].ProvinceName,   //string    `json:"province_name" gorm:"Column:province_name"`     // 省份
				CityName:       publishs[i].CityName,       //string    `json:"city_name" gorm:"Column:city_name"`
			})
		}

		fmt.Println(retPublish)

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": retPublish,
		})

	}

	fmt.Println("-----------------------------------------")
	fmt.Println("推荐列表结束")
	fmt.Println("-----------------------------------------")
}

func LikePublish(c *gin.Context) {
	fmt.Println(c)
	fmt.Println("喜欢内容")
	var reqUserLikePublish Entity.ReqUserLikePublish
	if c.BindJSON(&reqUserLikePublish) == nil { // 参数校验正确
		fmt.Println(reqUserLikePublish)

		redisUserInfo := Libaries.GetUserIdByToken(reqUserLikePublish.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]

		postId, _ := strconv.Atoi(string(reqUserLikePublish.PostId))
		userId, _ := strconv.Atoi(dbUserId.(string))

		if reqUserLikePublish.State == 1 {
			fmt.Println("新增")
			// 1. 加入喜爱
			userFocus := &Model.PublishUserLike{
				PublishId: postId,
				UserId:    userId,
				CreatedAt: time.Now(),
			}
			Db.DBHelper.Create(userFocus)

			// 2. 修改文章喜欢数量
			var publish Model.Publish
			Db.DBHelper.First(&publish, postId)
			publish.LikeCount += 1
			Db.DBHelper.Save(&publish)

		} else {
			fmt.Println("删除")
			// 1. 删除喜爱
			//var publishUserLike Model.PublishUserLike
			//Db.DBHelper.Where(&Model.PublishUserLike{PublishId: postId}).First(&publishUserLike)
			//Db.DBHelper.Debug().Delete(&publishUserLike)
			Db.DBHelper.Delete(Model.PublishUserLike{}, "publish_id = ?", postId)

			// 2. 修改文章喜欢数量
			var publish Model.Publish
			Db.DBHelper.First(&publish, postId)
			publish.LikeCount -= 1
			Db.DBHelper.Save(&publish)
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": "succ",
		})
	}

}

func PublishReply(c *gin.Context) {
	fmt.Println("请求回复列表")
	var reqPublishReply Entity.ReqPublishReply
	if c.BindJSON(&reqPublishReply) == nil { // 参数校验正确

		var publishReply []Model.PublishReply
		Db.DBHelper.Where("publish_id = ? AND status = 1", reqPublishReply.PublishId).Find(&publishReply)

		timeLayout := "2006-01-02 15:04:05" //转化所需模板
		var retPublishReply []Entity.RetPublishReply
		for i := 0; i < len(publishReply); i++ {
			retPublishReply = append(retPublishReply, Entity.RetPublishReply{
				Id:        publishReply[i].Id,                           //      int       `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"` // 主键编号
				PublishId: publishReply[i].PublishId,                    // int       `json:"publish_id" gorm:"Column:publish_id"`    // 内容编号
				UserId:    publishReply[i].UserId,                       //    int       `json:"user_id" gorm:"Column:user_id"`          // 用户编号
				NickName:  publishReply[i].NickName,                     //  string    `json:"nick_name" gorm:"Column:nick_name"`      // 昵称
				ReplyText: publishReply[i].ReplyText,                    // string    `json:"reply_text" gorm:"Column:reply_text"`    // 内容
				Status:    publishReply[i].Status,                       //    int       `json:"status" gorm:"Column:status"`            // 审核状态 0 未审核  1 审核
				CreatedAt: publishReply[i].CreatedAt.Format(timeLayout), // string `json:"created_at" gorm:"created_at"`           // 创建时间
			})
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": retPublishReply,
		})

	}
}

func Reply(c *gin.Context) {
	fmt.Println("回复")
	var reply Entity.Reply
	if c.BindJSON(&reply) == nil { // 参数校验正确

		redisUserInfo := Libaries.GetUserIdByToken(reply.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]
		nickName := redisUserInfo[1]

		userId, _ := strconv.Atoi(dbUserId.(string))
		dbNickName, _ := nickName.(string)
		reply := &Model.PublishReply{
			PublishId: reply.PublishId,
			UserId:    userId,
			NickName:  dbNickName,
			ReplyText: reply.ReplyText,
			Status:    0,

			CreatedAt: time.Now(),
		}
		Db.DBHelper.Create(reply)

	}
}
