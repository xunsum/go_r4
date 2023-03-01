package main

import (
	"github.com/gin-gonic/gin"
	"go_round4/routers"
)

func main() {
	r := gin.Default()
	//api v1 路由
	r.Group("/api/v1")

	//用户路由
	routers.UserRoutersInit(r)

	//视频内容相关路由
	routers.VideoRoutersInit(r)

	//admin
	routers.AdminRoutersInit(r)

	_ = r.Run(":9060")
}
