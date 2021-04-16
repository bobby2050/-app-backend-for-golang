package Middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("中间件检测token")

		//path := context.FullPath()
		//method := context.Request.Method
		//fmt.Println("path-->"+path)
		//fmt.Println("method--->"+method)

		//var token = context.Query("token")
		//token := ctx.Request.Header.Get("token")
		fmt.Println("--------")
		token := ctx.PostForm("token")

		if token == "" {
			ctx.JSON(http.StatusOK,  gin.H{"data":"","code": 900000, "msg":"token不存在"} )
			ctx.Abort()
		} else {
			// 得到token后，去redis中查询是否存在或者过期
			fmt.Println("token:", token)
		}

		ctx.Next()
	}
}
