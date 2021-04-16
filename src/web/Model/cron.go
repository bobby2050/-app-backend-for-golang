package Model

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"meizi/src/Db"
	"meizi/src/Redis"
	"meizi/src/web/CachEntity"
	"strconv"
)

// 发布审核
func PublishC() bool {
	var publish []Publish
	Db.DBHelper.Where("status", 0).Find(&publish)
	for i := 0; i < len(publish); i++ {
		//publish[i].PublishText 没有敏感字眼放行
		publishId := publish[i].Id
		if true {
			// 1.修改状态为成功并且保存到数据库
			publish[i].Status = 1
			Db.DBHelper.Save(&publish[i])
		} else {
			log.Println("[发布]审核失败，编号为: ", publishId)
		}

	}

	sql := "SELECT p.*, u.username, u.nick_name, u.email, u.city, p.created_at " +
		"FROM tbl_publish p " +
		"LEFT JOIN tbl_user u ON ( u.id = p.user_id  ) " +
		"WHERE p.status=1 ORDER BY p.id "

	var cachePublish []CachEntity.CachePublish
	Db.DBHelper.Raw(sql).Scan(&cachePublish)

	//timeLayout := "2006-01-02 15:04:05" //转化所需模板

	for j := 0; j < len(cachePublish); j++ {
		singlePostwriteToRedis(cachePublish[j])
	}

	return true
}

// 回复
func Reply() bool {
	var publishReply []PublishReply
	Db.DBHelper.Where("status", 0).Find(&publishReply)

	for i := 0; i < len(publishReply); i++ {
		//publishReply[i].ReplyText // 没有敏感字眼放行
		publishId := publishReply[i].PublishId // 贴子id
		if true {

			publishReply[i].Status = 1
			Db.DBHelper.Save(&publishReply[i])

			// 帖子数加1
			//var publish Publish
			//Db.DBHelper.Where("id = ?", publishId).Find(&publish)
			//publish.ReplyCount += 1
			//Db.DBHelper.Save(&publish)
			singlePublishReplyToRedis(publishReply[i])
		} else {
			log.Println("[回复]审核失败，编号为: ", publishId)
		}

	}
	//Db.DBHelper.Save(&publishReply)
	//Db.DBHelper.Create()

	return true
}

// 1.Post数据全部写入redis中
var ctx = context.Background()

func InitPostToRedis(c *gin.Context) {
	sql := "SELECT p.* , u.username, u.nick_name, u.email, u.city, p.created_at " +
		"FROM tbl_publish p " +
		"LEFT JOIN tbl_user u ON ( u.id = p.user_id  ) " +
		"WHERE p.status=1   ORDER BY p.id "

	var cachePublish []CachEntity.CachePublish
	Db.DBHelper.Raw(sql).Scan(&cachePublish)

	//timeLayout := "2006-01-02 15:04:05" //转化所需模板

	for i := 0; i < len(cachePublish); i++ {
		//cachePublish[i].CreatedAt.Format(timeLayout)
		singlePostwriteToRedis(cachePublish[i])
	}
}

// 单条写入redis
func singlePostwriteToRedis(post CachEntity.CachePublish) {
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	var retCachePublish CachEntity.RetCachePublish

	//DefaultTimeLoc := time.Local
	//loginTime, err := time.ParseInLocation("2006-01-02 15:04:05", lastLoginTime, DefaultTimeLoc)

	//fmt.Print(post.Id)
	//fmt.Print(post.CreatedAt)
	//fmt.Print("\r\n")

	retCachePublish = CachEntity.RetCachePublish{
		Id:             post.Id,
		UserId:         post.UserId,                       // 用户编号
		UserName:       post.UserName,                     // 用户名
		NickName:       post.NickName,                     // 昵称
		PublishText:    post.PublishText,                  // 发布内容
		Status:         post.Status,                       // 审核状态 0 未审核  1 审核
		CompanyName:    post.CompanyName,                  // 公司名称
		CreatedAt:      post.CreatedAt.Format(timeLayout), // 创建时间
		UpdatedAt:      post.UpdatedAt.Format(timeLayout), // 更新时间
		DepartmentName: post.DepartmentName,               // 部门
		GroupName:      post.GroupName,                    // 组名
		PositionName:   post.PositionName,                 // 职位
		ProvinceName:   post.ProvinceName,                 // 省份
		CityName:       post.CityName,                     // 城市
		//ReplyCount : post.ReplyCount,        // 回复数
		//LikeCount  : post.LikeCount,          // 喜欢数
	}
	responseStr, _ := json.Marshal(retCachePublish)
	s := string(responseStr)
	//fmt.Println("-----------------------------------------------------")
	//print(s)

	// ZSET
	member := &redis.Z{
		Score:  float64(post.Id),
		Member: s,
	}
	Redis.RedisPool.ZAdd(ctx, "zPostExample", member)
}

// 2.用户关注publish 数据全部写入redis中
func InitFocusToRedis(c *gin.Context) {
	sql := "SELECT f.* " +
		"FROM tbl_user_focus_publish f " +
		"ORDER BY f.id"

	var cacheFocus []CachEntity.CacheFocusPublish
	Db.DBHelper.Raw(sql).Scan(&cacheFocus)
	for i := 0; i < len(cacheFocus); i++ {
		singleFocusToRedis(cacheFocus[i])
	}
}

func singleFocusToRedis(row CachEntity.CacheFocusPublish) {
	Redis.RedisPool.SAdd(ctx, "sFocus:userid:"+strconv.Itoa(row.UserId), row.PublishId)
}

// 3.用户喜欢文章 数据全部写入redis中
func InitUserLovePublishToRedis(c *gin.Context) {
	sql := "SELECT l.publish_id, l.user_id " +
		"FROM tbl_publish_user_like l "
	var userLikePublish []CachEntity.UserLikePublish
	Db.DBHelper.Raw(sql).Scan(&userLikePublish)
	for i := 0; i < len(userLikePublish); i++ {
		singleUserLovePublishToRedis(userLikePublish[i])
	}
}

func singleUserLovePublishToRedis(row CachEntity.UserLikePublish) {
	Redis.RedisPool.SAdd(ctx, "sUserLikePublish:publishid:"+strconv.Itoa(row.PublishId), row.UserId)
}

// 4.文章回复 数据全部写入redis中
func InitPublishReplyToRedis(c *gin.Context) {
	sql := "SELECT * FROM tbl_publish_reply WHERE status = 1"
	var publishReply []PublishReply
	Db.DBHelper.Raw(sql).Scan(&publishReply)
	for i := 0; i < len(publishReply); i++ {
		singlePublishReplyToRedis(publishReply[i])
	}
}
func singlePublishReplyToRedis(row PublishReply) {
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	var retPublishReply CachEntity.RetPublishReply
	retPublishReply = CachEntity.RetPublishReply{
		Id:        row.Id,                           //      int       `json:"id" gorm:"PRIMARY_KEY" gorm:"Column:id"` // 主键编号
		PublishId: row.PublishId,                    //int       `json:"publish_id" gorm:"Column:publish_id"`    // 内容编号
		UserId:    row.UserId,                       // int       `json:"user_id" gorm:"Column:user_id"`          // 用户编号
		NickName:  row.NickName,                     //string    `json:"nick_name" gorm:"Column:nick_name"`      // 昵称
		ReplyText: row.ReplyText,                    // string    `json:"reply_text" gorm:"Column:reply_text"`    // 内容
		Status:    row.Status,                       //  int       `json:"status" gorm:"Column:status"`            // 审核状态 0 未审核  1 审核
		CreatedAt: row.CreatedAt.Format(timeLayout), //string `json:"created_at" gorm:"created_at"`           // 创建时间
	}
	en, _ := json.Marshal(retPublishReply)

	// ZSET
	member := &redis.Z{
		Score:  float64(row.Id),
		Member: en,
	}
	Redis.RedisPool.ZAdd(ctx, "zPublishReplyExample:publishId:"+strconv.Itoa(row.PublishId), member)
}

// keys zPublishReplyExample*
