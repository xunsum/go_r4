package routers

import (
	"github.com/gin-gonic/gin"
	"go_round4/controllers/adminControllers"
	"go_round4/middlewares"
	"go_round4/res"
)

func AdminRoutersInit(r *gin.Engine) {

	adminRouters := r.Group("/admin")
	{
		adminRouters.GET("/login", adminControllers.Login)
		adminRouters.GET("/encPas", adminControllers.EncryptPassword)
	}
	adminRouters.Use(middlewares.AuthJWT(res.ADMIN_MODE))
	{
		adminRouters.PUT("/setUserStatus", adminControllers.SetUserStatus)
		adminRouters.PUT("/setCommentVisibility", adminControllers.SetCommentVisibility)
	}
}
