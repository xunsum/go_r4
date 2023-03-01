package routers

import (
	"github.com/gin-gonic/gin"
	"go_round4/controllers/videoControllers"
	"go_round4/middlewares"
	"go_round4/res"
)

func VideoRoutersInit(r *gin.Engine) {
	videoRouters := r.Group("/video")
	videoRouters.Use(middlewares.AuthJWT(res.USER_MODE), middlewares.AuthId)
	videoRouters.POST("/upload", videoControllers.UploadVideo)
	videoRouters.PUT("/setLike", videoControllers.SetLikeVideo)
	videoRouters.PUT("/setComment", videoControllers.SetComment)
	videoRouters.PUT("/setCollect", videoControllers.SetCollection)
	videoRouters.GET("/searchVideo", videoControllers.SearchVideo)
	videoRouters.PUT("/putDanmaku", videoControllers.SetDanmaku)
	videoRouters.GET("/getVideo", videoControllers.GetVideo)
	videoRouters.GET("/getDanmakuList", videoControllers.GetDanmakuList)
	videoRouters.GET("/getVideoInfo", videoControllers.GetVideoInfo)
	videoRouters.GET("/getComments", videoControllers.GetComments)
}
