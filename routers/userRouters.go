package routers

import (
	"github.com/gin-gonic/gin"
	"go_round4/controllers/userControllers"
	"go_round4/middlewares"
)

func UserRoutersInit(r *gin.Engine) {
	unAccessedUserRouters := r.Group("/welcome")
	{
		//注册
		unAccessedUserRouters.POST("/register", userControllers.Register)

		//登录
		unAccessedUserRouters.GET("/loginViaMail", userControllers.Login(userControllers.NameLogin))
		unAccessedUserRouters.GET("/loginViaUserName", userControllers.Login(userControllers.MailLogin))

	}

	userRouters := r.Group("/user")
	userRouters.Use(middlewares.AuthJWT(0), middlewares.AuthId)
	{
		userRouters.PUT("/putSlogan", userControllers.PutSlogan)
		userRouters.GET("/getUserHead", userControllers.GetAvatar)
		userRouters.GET("/search", userControllers.SearchUser)
		userRouters.POST("/putUserHead", userControllers.PostUserHead)
		userRouters.PUT("/putUserEmail", userControllers.PutUserEmail)
		userRouters.PUT("/putPassword", userControllers.SetPassword)
		userRouters.GET("/getUserInfo", userControllers.GetUserInfo)
		userRouters.GET("/getCollections", userControllers.GetCollections)
		userRouters.GET("/getSearchHistory", userControllers.GetSearchHistories)
	}
}
