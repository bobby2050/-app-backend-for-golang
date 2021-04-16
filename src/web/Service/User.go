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

var redisCtx = context.Background()

// 注册
func Register(c *gin.Context)  {
	fmt.Println("注册")
	var reqUserRegistion Entity.ReqUserRegistion

	if c.BindJSON(&reqUserRegistion) == nil {
		fmt.Println("从客户端接收到的：")
		fmt.Println("UserName:", reqUserRegistion.UserName)
		fmt.Println("Password:", reqUserRegistion.Password)
		fmt.Println("RePassword:", reqUserRegistion.RePassword)
		fmt.Println("NickName:", reqUserRegistion.NickName)
		fmt.Println("Email:", reqUserRegistion.Email)


		var salt = Libaries.CreateRandomString(4)
		var password = Libaries.Md5V(salt + reqUserRegistion.Password)
		user := &Model.User{
			UserName: reqUserRegistion.UserName,
			Password: password,
			Salt : salt,
			NickName: reqUserRegistion.NickName,
			Email: reqUserRegistion.Email,
			CreatedAt: time.Now(),
		}
		//Db.DBHelper.Create(user).Error
		Db.DBHelper.Create(user)

		fmt.Println("id: " + strconv.Itoa(user.Id))

		if user.Id == 0 {
			c.JSON(200, gin.H{
				"code": 90000,
				"msg": "fail",
				"data": "",

			})
		} else {
			c.JSON(200, gin.H{
				"code": 10000,
				"msg": "success",
				"data": "",

			})
		}

	}

}

// 登陆
func Login(c *gin.Context)  {


	//db.Where("typeId=?", "abc").Find(&tc)

	//Db.DBHelper.Where(&Model.User{UserName:"abc"}).Find(&tc)

	var userlogin Entity.ReqUserLogin
	if c.BindJSON(&userlogin) == nil {
		fmt.Println("从客户端接收到的：")
		fmt.Println("Name:", userlogin.Name)
		fmt.Println("Password:", userlogin.Password)

		var user Model.User
		Db.DBHelper.Where(&Model.User{UserName: userlogin.Name}).First(&user)
		//fmt.Println("------------------------")
		//fmt.Println(user)
		//fmt.Println(user.Salt)
		//fmt.Println(userlogin.Name)
		//fmt.Println("------------------------")
		var loginUserData Entity.LoginUserData
		if user.Password == Libaries.Md5V(user.Salt + userlogin.Password) {
			fmt.Println("密码匹配成功")
			// 计算token值
			userToken := Libaries.Md5V( strconv.Itoa(user.Id ) + string(time.Now().UnixNano() ) )


			loginUserData.Token = userToken
			loginUserData.UserInfo = user
			fmt.Println(user)

			// 保存用户信息到redis
			err := Redis.RedisPool.HMSet(ctx,  "token_" + userToken, "username", user.UserName, "nick_name", user.NickName, "user_id", user.Id, "mobile", user.Mobile).Err()
			if err != nil {
				panic(err)
			}

			c.JSON(200, gin.H{
				"msg": "success",
				"data": loginUserData,
				"code": 10000,
			})
		} else {
			fmt.Println("密码匹配失败")
			c.JSON(200, gin.H{
				"msg": "fail",
				"data": loginUserData,
				"code": 90000,
			})
		}



		//c.JSON(200, gin.H{
		//	"msg": "success",
		//	"data": "",
		//	"code": 100000,
		//})
	}

}

func FocusPublish(c *gin.Context)  {

	fmt.Println(c)
	fmt.Println("关注用户")
	var reqFocusPublish Entity.ReqFocusPublish
	if c.BindJSON(&reqFocusPublish) == nil { // 参数校验正确
		fmt.Println(reqFocusPublish)

		redisUserInfo := Libaries.GetUserIdByToken(reqFocusPublish.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]

		publishId :=  reqFocusPublish.PublishId
		userId, _ := strconv.Atoi(dbUserId.(string))
		userFocus := &Model.UserFocusPublish{
			UserId:      userId,
			PublishId: publishId,
			CreatedAt:   time.Now(),
		}
		//Db.DBHelper.Create(user).Error
		Db.DBHelper.Create(userFocus)
		c.JSON(200, gin.H{
			"code": 10000,
			"msg": "success",
			"data": "",

		})

	}
}

// 取消关注用户
func CancelFocusUser(c *gin.Context) {
	fmt.Println("取消关注用户")
	var reqFocusPublish Entity.ReqFocusPublish
	if c.BindJSON(&reqFocusPublish) == nil { // 参数校验正确
		fmt.Println(reqFocusPublish)

		redisUserInfo := Libaries.GetUserIdByToken(reqFocusPublish.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]

		tid :=  reqFocusPublish.PublishId // 这里接收到的是表的id，并非userid
		userId, _ := strconv.Atoi(dbUserId.(string))

		//Db.DBHelper.Create(user).Error
		Db.DBHelper.Delete(Model.UserFocusPublish{}, "id = ?  AND user_id = ?", tid, userId)
		c.JSON(200, gin.H{
			"code": 10000,
			"msg": "success",
			"data": "",

		})

	}
}





// 我的关注
func MyFocus(c *gin.Context) {
	var myFocusUser Entity.MyFocusUser

	if c.BindJSON(&myFocusUser) == nil {
		fmt.Println("Token:" + myFocusUser.Token)
		redisUserInfo := Libaries.GetUserIdByToken(myFocusUser.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]
		userId, _ := strconv.Atoi(dbUserId.(string))

		var publishs []Entity.MyFocus

		sql := "SELECT f.id, f.created_at, u.nick_name, u.city, f.user_id " +
			"FROM tbl_user_focus f " +
			"LEFT JOIN tbl_user u on u.id = f.focus_user_id " +
			"WHERE f.user_id = ? " +
			"ORDER BY f.created_at DESC"

		Db.DBHelper.Raw(sql, userId).Scan(&publishs)

		timeLayout := "2006-01-02 15:04:05" //转化所需模板
		var retMyFocus []Entity.RetMyFocus
		for i:= 0; i< len(publishs) ; i++ {
			retMyFocus = append(retMyFocus,Entity.RetMyFocus{
				Id  : publishs[i].Id,
				UserId: publishs[i].UserId,
				CreatedAt : publishs[i].CreatedAt.Format(timeLayout),
				NickName : publishs[i].NickName,
				City: publishs[i].City,
			} )
		}
		c.JSON(200, gin.H{
			"code": 10000,
			"msg": "success",
			"data": retMyFocus,

		})


	}
}

// 我的发布
func MyPublish(c *gin.Context) {
	var myPublish Entity.MyPublish
	if c.BindJSON(&myPublish) == nil {

		redisUserInfo := Libaries.GetUserIdByToken(myPublish.Token)
		fmt.Println(redisUserInfo)
		dbUserId := redisUserInfo[0]
		userId, _ := strconv.Atoi(dbUserId.(string))

		var publish []Model.Publish
		Db.DBHelper.Where("user_id = ?", userId).Find(&publish)



		timeLayout := "2006-01-02 15:04:05" //转化所需模板
		var retPublish []Entity.RetPublish
		for i := 0; i < len(publish); i++ {
			retPublish = append(retPublish, Entity.RetPublish{
				Id:             publish[i].Id,
				UserId:         publish[i].UserId,
				UserName:       publish[i].UserName,
				NickName:       publish[i].NickName,
				PublishText:    publish[i].PublishText,
				Status:         publish[i].Status,
				//ReplyCount:     publish[i].ReplyCount,
				//LikeCount:      publish[i].LikeCount,
				CompanyName:    publish[i].CompanyName,
				CreatedAt:      publish[i].CreatedAt.Format(timeLayout),
				UpdatedAt:      publish[i].UpdatedAt.Format(timeLayout),
				DepartmentName: publish[i].DepartmentName, //
				GroupName:      publish[i].GroupName,      //string    `json:"group_name" gorm:"Column:group_name"`           // 组名
				PositionName:   publish[i].PositionName,   //string    `json:"position_name" gorm:"Column:position_name"`     // 职位
				ProvinceName:   publish[i].ProvinceName,   //string    `json:"province_name" gorm:"Column:province_name"`     // 省份
				CityName:       publish[i].CityName,       //string    `json:"city_name" gorm:"Column:city_name"`             // 城市
			})
		}




		c.JSON(200, gin.H{
			"code": 10000,
			"msg": "success",
			"data": retPublish,

		})

	}
}