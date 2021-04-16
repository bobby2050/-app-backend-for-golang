package main

import (
	"github.com/gin-gonic/gin"
	"meizi/src/web/Router"
)




func main() {

	r := gin.New()
	Router.Route(r)
	r.Run()

}
