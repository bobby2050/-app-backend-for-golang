package Cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"meizi/src/Db"
	"meizi/src/web/CachEntity"
	"meizi/src/web/Model"
	"strconv"
	"time"

	//"github.com/go-redis/redis/v8"

	"meizi/src/Redis"
	"meizi/src/web/Entity"
	"meizi/src/web/Libaries"
)

// 举报
func ReprotPost(c *gin.Context)  {

	var reportPost Entity.ReportPost
	if c.BindJSON(&reportPost) == nil {
		var userId string
		if reportPost.Token != "" {
			redisUserInfo := Libaries.GetUserIdByToken(reportPost.Token)

			dbUserId := redisUserInfo[0]
			userId = dbUserId.(string)
		}

		if reportPost.ReportType == 1 {
			fmt.Println("userId: " + userId + ",类型:色情，Post：" + strconv.Itoa(int(reportPost.PostId)))

		} else if reportPost.ReportType == 2{
			fmt.Println("userId: " + userId + ",类型:政治敏感，Post：" + strconv.Itoa(int(reportPost.PostId)))

		} else if reportPost.ReportType == 3{
			fmt.Println("userId: " + userId + ",类型:广告，Post：" + strconv.Itoa(int(reportPost.PostId)))

		} else if reportPost.ReportType == 4{
			fmt.Println("userId: " + userId + ",类型:违纪违法，Post：" + strconv.Itoa(int(reportPost.PostId)))
			
		}

	}
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": "",
	})
}
func MyPublish(c *gin.Context)  {
	fmt.Println("我的发布列表")
	var ctx = context.Background()
	var myPublish Entity.MyPublish
	if c.BindJSON(&myPublish) == nil {
		redisUserInfo := Libaries.GetUserIdByToken(myPublish.Token)
		//fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]

		//fmt.Println(dbUserId)


		userIdString, _ := dbUserId.(string)


		//fmt.Println(userIdString)
		//fmt.Println("---------")
		userId, _ := strconv.Atoi(userIdString)
		//fmt.Println(userId)

		redisRange := &redis.ZRangeBy{
			Max:"+inf",
			Min:"-inf",
			Offset:0,
			Count: -1,
		}
		res, err := Redis.RedisPool.ZRevRangeByScore(ctx,"zPostExample", redisRange).Result()
		if err != nil {
			panic(err)
		}

		myMap := make(map[int]CachEntity.RetCachePublish, 5)

		//fmt.Println("============================")
		for i := 0 ; i < len(res); i++ {
			var cachePublish CachEntity.RetCachePublish

			json.Unmarshal([]byte(res[i]), &cachePublish)

			if cachePublish.UserId == userId {

				//myMap = append(myMap, retCachePublishs)
				//myMap =  append(myMap, cachePublish)
				myMap[cachePublish.Id] = cachePublish
			}
		}

		var retCachePublishs []CachEntity.RetCachePublish
		for j := 0 ;j < len(myMap); j++ {
			if myMap[j].Id > 0 {
				//fmt.Println("------------------------------------")
				//fmt.Println(myMap[j])
				retCachePublishs = append(retCachePublishs, myMap[j])
			}
		}

		//fmt.Println("============================")
		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": retCachePublishs,
		})

	}
}

func CancelFocusPublish(c *gin.Context)  {
	var ctx = context.Background()
	fmt.Println("取消关注用户")
	var reqFocusPublish Entity.ReqFocusPublish
	if c.BindJSON(&reqFocusPublish) == nil { // 参数校验正确
		fmt.Println(reqFocusPublish)

		redisUserInfo := Libaries.GetUserIdByToken(reqFocusPublish.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]
		u := dbUserId.(string)

		Redis.RedisPool.SRem(ctx, "sFocus:userid:" + u, reqFocusPublish.PublishId)

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": "",
		})
	}
}
func FocusPublish(c *gin.Context)  {
	var ctx = context.Background()
	fmt.Println("关注用户")
	var reqFocusPublish Entity.ReqFocusPublish
	if c.BindJSON(&reqFocusPublish) == nil { // 参数校验正确
		fmt.Println(reqFocusPublish)

		redisUserInfo := Libaries.GetUserIdByToken(reqFocusPublish.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]

		// 1.写入redis
		u := dbUserId.(string)
		Redis.RedisPool.SAdd(ctx, "sFocus:userid:" + u, reqFocusPublish.PublishId)

		// 2.写入mysql
		//userFocus := &Model.UserFocusPublish{
		//	UserId:      dbUserId.(int),
		//	PublishId: reqFocusPublish.PublishId,
		//	CreatedAt:   time.Now(),
		//}
		//Db.DBHelper.Create(userFocus)

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": "",
		})
	}
}

func PostReplyLikeCount(c *gin.Context) {
	var ctx = context.Background()
	fmt.Println("回复和喜欢数")
	var postReplyLike Entity.PostReplyLike
	if c.BindJSON(&postReplyLike) == nil { // 参数校验正确

		redisUserInfo := Libaries.GetUserIdByToken(postReplyLike.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]

		var uid string
		if dbUserId != nil {
			uid = dbUserId.(string)
		}

		var retPostReplyLike []CachEntity.RetPostReplyLike
		for i:= 0 ; i< len(postReplyLike.PostId); i++ {

			redisRange := &redis.ZRangeBy{
				Max:"+inf",
				Min:"-inf",
				Offset:0,
				Count: -1,
			}
			replyRes, err := Redis.RedisPool.ZRevRangeByScore(ctx,"zPublishReplyExample:publishId:" + strconv.Itoa(int(postReplyLike.PostId[i])), redisRange).Result()

			var replyCount int
			for i := 0; i < len(replyRes); i++ {
				var retPublishReply CachEntity.RetPublishReply
				json.Unmarshal([]byte(replyRes[i]), &retPublishReply)
				if retPublishReply.Status == 1 {
					replyCount += 1
				}
			}


			//replyCount, err := Redis.RedisPool.ZCard(ctx, "zPublishReplyExample:publishId:" + strconv.Itoa(int(postReplyLike.PostId[i]))  ).Result()
			//if err != nil {
			//	panic(err)
			//}

			likeCount, err := Redis.RedisPool.SCard(ctx, "sUserLikePublish:publishid:" + strconv.Itoa(int(postReplyLike.PostId[i]))  ).Result()
			if err != nil {
				panic(err)
			}

			isFocus, err := Redis.RedisPool.SIsMember(ctx, "sFocus:userid:" + uid ,postReplyLike.PostId[i] ).Result()
			if err != nil {
				panic(err)
			}

			fmt.Println(postReplyLike.PostId[i])

			//retPostReplyLike[i].PostId = postReplyLike.PostId[i]
			var row CachEntity.RetPostReplyLike
			row.PostId = postReplyLike.PostId[i]
			row.PostReplyLike.LikeNum = int(likeCount)
			row.PostReplyLike.ReplyNum = int(replyCount)
			row.PostReplyLike.IsFocus = !isFocus

			retPostReplyLike = append(retPostReplyLike, row)
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": retPostReplyLike,
		})

	}
}
func Reply(c *gin.Context) {
	//var ctx = context.Background()
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




		//timeLayout := "2006-01-02 15:04:05" //转化所需模板
		//var retPublishReply CachEntity.RetPublishReply
		//retPublishReply = CachEntity.RetPublishReply{
		//	Id  : reply.Id, //      int       `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"` // 主键编号
		//	PublishId :reply.PublishId, //int       `json:"publish_id" gorm:"Column:publish_id"`    // 内容编号
		//	UserId   : reply.UserId, // int       `json:"user_id" gorm:"Column:user_id"`          // 用户编号
		//	NickName  : reply.NickName, //string    `json:"nick_name" gorm:"Column:nick_name"`      // 昵称
		//	ReplyText: reply.ReplyText, // string    `json:"reply_text" gorm:"Column:reply_text"`    // 内容
		//	Status  : reply.Status, //  int       `json:"status" gorm:"Column:status"`            // 审核状态 0 未审核  1 审核
		//	CreatedAt : reply.CreatedAt.Format(timeLayout), //string `json:"created_at" gorm:"created_at"`           // 创建时间
		//}
		//en,_ := json.Marshal(retPublishReply)
		//
		//// ZSET
		//member := &redis.Z{
		//	Score: float64(reply.Id),
		//	Member: en,
		//}
		//Redis.RedisPool.ZAdd(ctx,"zPublishReplyExample:publishId:" + strconv.Itoa(reply.PublishId), member)

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": "",
		})

	}
}

func LikePublish(c *gin.Context )  {
	var ctx = context.Background()
	fmt.Println(c)
	fmt.Println("喜欢内容")
	var reqUserLikePublish Entity.ReqUserLikePublish
	if c.BindJSON(&reqUserLikePublish) == nil { // 参数校验正确
		fmt.Println(reqUserLikePublish)

		redisUserInfo := Libaries.GetUserIdByToken(reqUserLikePublish.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]


		userId, _ := strconv.Atoi(dbUserId.(string))

		if reqUserLikePublish.State == 1 {

			Redis.RedisPool.SAdd(ctx, "sUserLikePublish:publishid:" + strconv.Itoa(int(reqUserLikePublish.PostId)), userId)
		} else {
			Redis.RedisPool.SRem(ctx, "sUserLikePublish:publishid:" + strconv.Itoa(int(reqUserLikePublish.PostId)), userId)
		}
		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": "",
		})
	}
}
// 文章回复列表
func PublishReply(c *gin.Context) {
	var ctx = context.Background()
	fmt.Println("-----------------------------------------")
	fmt.Println("回复 开始")
	fmt.Println("-----------------------------------------")
	var reqPublishReplyList Entity.ReqPublishReplyList
	if c.Bind(&reqPublishReplyList) == nil { // 参数校验正确

		redisRange := &redis.ZRangeBy{
			Max:    "+inf",
			Min:    "-inf",
			//Offset: 1,
			//Count:  5,
		}
		res, err := Redis.RedisPool.ZRangeByScore(ctx, "zPublishReplyExample:publishId:" + strconv.Itoa(reqPublishReplyList.PublishId), redisRange).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(res)

		var retPublishReply []CachEntity.RetPublishReply
		for i:=0;i<len(res);i++ {
			var reply CachEntity.RetPublishReply
			json.Unmarshal([]byte(res[i]), &reply)

			if reply.Status == 1{
 				retPublishReply = append(retPublishReply, reply)
			}
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": retPublishReply,
		})


	}
}

// 关注列表
func FocusPublishList(c *gin.Context) {
	var ctx = context.Background()
	fmt.Println("-----------------------------------------")
	fmt.Println("列表列表 开始")
	fmt.Println("-----------------------------------------")
	var reqUserFocusPublishList Entity.ReqUserFocusPublishList
	if c.Bind(&reqUserFocusPublishList) == nil { // 参数校验正确
		token := reqUserFocusPublishList.Token // 获取Get参数 token
		fmt.Println("token： " + token)
		redisUserInfo := Libaries.GetUserIdByToken(token)
		fmt.Println("key", redisUserInfo)
		dbUserId := redisUserInfo[0]

		userId := dbUserId.(string)
		res, err := Redis.RedisPool.SMembers(ctx, "sFocus:userid:" + userId).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(res)

		var fRetPublish []Entity.RetPublish
		for i := 0 ; i< len(res); i++ {
			//fmt.Println(res[i])
			id,_ := strconv.Atoi(res[i])
			retCachePublishs := getDataByid(ctx,  id)
			var retPublish Entity.RetPublish
			json.Unmarshal([]byte(retCachePublishs), &retPublish)

			likeCount, err := Redis.RedisPool.SCard(ctx, "sUserLikePublish:publishid:" + strconv.Itoa(retPublish.Id)  ).Result()
			if err != nil {
				panic(err)
			}

			replyCount, err := Redis.RedisPool.ZCard(ctx, "zPublishReplyExample:publishId:" + strconv.Itoa(retPublish.Id)  ).Result()
			if err != nil {
				panic(err)
			}

			//fmt.Print(res)

			fRetPublish = append(fRetPublish, Entity.RetPublish{
				Id:             retPublish.Id,
				UserId:         retPublish.UserId,
				UserName:       retPublish.UserName,
				NickName:       retPublish.NickName,
				PublishText:    retPublish.PublishText,
				Status:         retPublish.Status,
				ReplyCount:     replyCount,
				LikeCount:      likeCount,
				CompanyName:    retPublish.CompanyName,
				CreatedAt:      retPublish.CreatedAt,
				UpdatedAt:      retPublish.UpdatedAt,
				//IsLike:         getUserLikePublish(retCachePublishs[i].Id, db_id),
				//IsFocus:        getUserFocusUser(retCachePublishs[i].UserId,  db_id ), //  用户编号
				DepartmentName: retPublish.DepartmentName, //
				GroupName:      retPublish.GroupName,      //string    `json:"group_name" gorm:"Column:group_name"`           // 组名
				PositionName:   retPublish.PositionName,   //string    `json:"position_name" gorm:"Column:position_name"`     // 职位
				ProvinceName:   retPublish.ProvinceName,   //string    `json:"province_name" gorm:"Column:province_name"`     // 省份
				CityName:       retPublish.CityName,       //string    `json:"city_name" gorm:"Column:city_name"`
			})


			//fRetPublish = append(fRetPublish, retPublish)
		}
		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": fRetPublish,
		})
	}

}

func getDataByid(ctx context.Context ,id int)  string {

	redisRange := &redis.ZRangeBy{
		Max:    strconv.Itoa(id),
		Min:    strconv.Itoa(id),
		//Offset: 1,
		//Count:  5,
	}
	res, err := Redis.RedisPool.ZRangeByScore(ctx, "zPostExample", redisRange).Result()
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	return res[0]
}

func MyfocusPostId(c *gin.Context)  {
	var ctx = context.Background()

	var reqToken Entity.ReqToken
	if c.Bind(&reqToken) == nil { // 参数校验正确
		token := reqToken.Token // 获取Get参数 token
		fmt.Println("token： " + token)
		redisUserInfo := Libaries.GetUserIdByToken(token)
		fmt.Println("key", redisUserInfo)
		dbUserId := redisUserInfo[0]
		//nick_name := redisUserInfo[1]
		//username := redisUserInfo[2]
		fmt.Println("user_id：", dbUserId)

		uid := dbUserId.(string)
		// sFocus:userid:20
		res, err := Redis.RedisPool.SMembers(ctx, "sFocus:userid:"+uid).Result()
		if err != nil {
			panic(err)
		}
		var ids [] int
		for i := 0; i < len(res); i++ {
			id, _ := strconv.Atoi(res[i])
			ids = append(ids,  id )
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": ids,
		})
	}
}

func MyPublishLike(c *gin.Context)  {
	var ctx = context.Background()

	var reqToken Entity.ReqToken
	if c.Bind(&reqToken) == nil { // 参数校验正确
		token := reqToken.Token // 获取Get参数 token
		fmt.Println("token： " + token)
		redisUserInfo := Libaries.GetUserIdByToken(token)
		fmt.Println("key", redisUserInfo)
		db_user_id := redisUserInfo[0]
		//nick_name := redisUserInfo[1]
		//username := redisUserInfo[2]
		fmt.Println("user_id：", db_user_id)


		redisRange := &redis.ZRangeBy{
			Max:"+inf",
			Min:"-inf",

		}
		res, err := Redis.RedisPool.ZRevRangeByScore(ctx,"zPostExample", redisRange).Result()
		if err != nil {
			panic(err)
		}

		var mylist1 []int
		for i := 0 ;i < len(res); i++ {
			var cachePublish CachEntity.RetCachePublish
			json.Unmarshal([]byte(res[i]), &cachePublish)
			fmt.Println()
			IsM, _ := Redis.RedisPool.SIsMember(ctx, "sUserLikePublish:publishid:"+strconv.Itoa(cachePublish.Id), db_user_id).Result()
			fmt.Println("cachePublish.Id:" + strconv.Itoa(cachePublish.Id) + "IsM:" )
			fmt.Println(IsM)
			fmt.Println("----")
			if IsM == true{
				mylist1 = append(mylist1, cachePublish.Id)
			}
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": mylist1,
		})


	}
}

// 推荐列表
func PublishList(c *gin.Context) {
	var ctx = context.Background()
	fmt.Println("-----------------------------------------")
	fmt.Println("推荐列表 开始")
	fmt.Println("-----------------------------------------")

	var reqUserFocusPublishList Entity.ReqUserFocusPublishList
	if c.Bind(&reqUserFocusPublishList) == nil { // 参数校验正确

		token := reqUserFocusPublishList.Token // 获取Get参数 token
		fmt.Println("token： " + token)
		redisUserInfo := Libaries.GetUserIdByToken(token)
		fmt.Println("key", redisUserInfo)
		db_user_id := redisUserInfo[0]
		//nick_name := redisUserInfo[1]
		//username := redisUserInfo[2]
		fmt.Println("user_id：", db_user_id)


		flag := reqUserFocusPublishList.Flag // 获取Get参数 page
		fmt.Println("flag: " + flag)

		reqPublishId := reqUserFocusPublishList.PublishId

		fmt.Print("req_publish_id:" + strconv.Itoa(reqPublishId))

		var res []string
		var err error
		if flag == "first" {
			fmt.Print("if first")
			redisRange := &redis.ZRangeBy{
				Max:    "+inf",
				Min:    strconv.Itoa(reqPublishId),
				Offset: 1,
				Count:  5,
			}
			res, err = Redis.RedisPool.ZRangeByScore(ctx, "zPostExample", redisRange).Result()
			if err != nil {
				panic(err)
			}
		} else if flag == "last"{
			fmt.Print("if last")
			redisRange := &redis.ZRangeBy{
				Max:    strconv.Itoa(reqPublishId),
				Min:    "-inf",
				Offset: 1,
				Count:  5,
			}
			res, err = Redis.RedisPool.ZRevRangeByScore(ctx, "zPostExample", redisRange).Result()
			if err != nil {
				panic(err)
			}

		} else {
			fmt.Print("if else")
			redisRange := &redis.ZRangeBy{
				Max:"+inf",
				Min:"-inf",
				Offset:0,
				Count: 5,
			}
			res, err = Redis.RedisPool.ZRevRangeByScore(ctx,"zPostExample", redisRange).Result()
			if err != nil {
				panic(err)
			}
		}




		fmt.Println(res)

		var retCachePublishs []CachEntity.RetCachePublish

		for i := 0 ;i < len(res); i++ {
			var cachePublish CachEntity.RetCachePublish
			json.Unmarshal([]byte(res[i]), &cachePublish)
			retCachePublishs = append(retCachePublishs, cachePublish)
		}

		// 返回结果
		var retPublish []Entity.RetPublish
		for i := 0; i < len(retCachePublishs); i++ {
			//db_id, _ := strconv.Atoi(fmt.Sprintf("%v", db_user_id))

			likeCount, err := Redis.RedisPool.SCard(ctx, "sUserLikePublish:publishid:" + strconv.Itoa(retCachePublishs[i].Id)  ).Result()
			if err != nil {
				panic(err)
			}

			replyCount, err := Redis.RedisPool.ZCard(ctx, "zPublishReplyExample:publishId:" + strconv.Itoa(retCachePublishs[i].Id)  ).Result()
			if err != nil {
				panic(err)
			}

			if checkRepeatData(retPublish, retCachePublishs[i].Id) {

				retPublish = append(retPublish, Entity.RetPublish{
					Id:             retCachePublishs[i].Id,
					UserId:         retCachePublishs[i].UserId,
					UserName:       retCachePublishs[i].UserName,
					NickName:       retCachePublishs[i].NickName,
					PublishText:    retCachePublishs[i].PublishText,
					Status:         retCachePublishs[i].Status,
					ReplyCount:     replyCount,
					LikeCount:      likeCount,
					CompanyName:    retCachePublishs[i].CompanyName,
					CreatedAt:      retCachePublishs[i].CreatedAt,
					UpdatedAt:      retCachePublishs[i].UpdatedAt,
					//IsLike:         getUserLikePublish(retCachePublishs[i].Id, db_id),
					//IsFocus:        getUserFocusUser(retCachePublishs[i].UserId,  db_id ), //  用户编号
					DepartmentName: retCachePublishs[i].DepartmentName, //
					GroupName:      retCachePublishs[i].GroupName,      //string    `json:"group_name" gorm:"Column:group_name"`           // 组名
					PositionName:   retCachePublishs[i].PositionName,   //string    `json:"position_name" gorm:"Column:position_name"`     // 职位
					ProvinceName:   retCachePublishs[i].ProvinceName,   //string    `json:"province_name" gorm:"Column:province_name"`     // 省份
					CityName:       retCachePublishs[i].CityName,       //string    `json:"city_name" gorm:"Column:city_name"`
				})
			}
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": retPublish,
		})

	}

	fmt.Println("-----------------------------------------")
	fmt.Println("推荐列表 结束")
	fmt.Println("-----------------------------------------")
}

func checkRepeatData(rows []Entity.RetPublish,id int) bool {
	flag := true
	for i := 0; i< len(rows); i++ {
		if rows[i].Id == id {
			flag = false
		}
	}
	return flag
}

// 用户关注用户
//func getUserFocusUser(userId int, db_user_id int) bool  {
//	var ctx = context.Background()
//	IsM, _ := Redis.RedisPool.SIsMember(ctx, "sFocus:userid:"+strconv.Itoa(userId), db_user_id).Result()
//	return IsM
//}

// 用户喜欢文章
//func getUserLikePublish(id int, userId int) bool  {
//	var ctx = context.Background()
//	fmt.Println("publishid:"+strconv.Itoa(id)+ ", userId:" + strconv.Itoa(userId))
//	IsM, _ := Redis.RedisPool.SIsMember(ctx, "sUserLikePublish:publishid:"+strconv.Itoa(id), userId).Result()
//	return IsM
//}
