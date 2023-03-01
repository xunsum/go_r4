package adminControllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_round4/models"
	"go_round4/models/response"
	"go_round4/res"
	"go_round4/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func Login(c *gin.Context) {
	var result *gorm.DB
	name := c.Query("user_name")
	password := c.Query("password")

	//搜索用户
	var searchOutcome models.Admin
	result = utils.DB.Where("name = ?", name).First(&searchOutcome)

	if fmt.Sprintf("%v", result.Error) == "record not found" {
		//无此用户：
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "No such user.",
			Error:  fmt.Sprintf("No such user. name: %s", name),
		})
		log.Printf("No such user. name: %s.", name)
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem searching user in databases.",
			Error:  fmt.Sprintf("User search error: %v, name: %s", result.Error, name),
		})
		log.Printf("User search error: %v, name: %s", result.Error, name)
		c.Abort()
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(searchOutcome.EncryptedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		//错误密码：
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "Wrong password",
			Error:  "Wrong password",
		})
		log.Printf("Wrong password: user: %v", name)
		c.Abort()
		return
	}

	if err == nil {
		// 签发 token
		token, err := utils.Sign(res.ADMIN_MODE, searchOutcome.Uid, searchOutcome.Name)
		if err != nil {
			utils.ShowUnknownTokenError(c, err, name, 1)
			c.Abort()
			return
		}

		c.JSON(200, response.AdminResponse{
			Status: 200,
			Data:   searchOutcome,
			Msg:    fmt.Sprintf("Welcome, admin: %v, your token: %v", name, token),
			Error:  fmt.Sprintf("token: %v", token),
		})
	}
}

func EncryptPassword(c *gin.Context) { // 摆烂写的，用 poster 获取一个加过密的密码，不想做管理员注册
	pass := c.Query("pass")
	encryptedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   string(encryptedPasswordBytes),
		Msg:    string(encryptedPasswordBytes),
		Error:  string(encryptedPasswordBytes),
	})
	c.Next()
}

func SetUserStatus(c *gin.Context) {
	targetUid := c.Query("uid")
	isBlockedString := c.Query("is_blocked")

	isBlocked, err := strconv.Atoi(isBlockedString)
	result := utils.DB.Where("uid = ?", targetUid).Model(models.User{}).Update("is_blocked", isBlocked)
	if err != nil || result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem block/unblocking user.",
			Error:  fmt.Sprintf("Having problem block/unblocking user, error: %v， %v", err, result.Error),
		})
		log.Printf("Having problem block/unblocking user, error: %v， %v", err, result.Error)
		c.Abort()
		return
	}
	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   "block/unblock user success.",
		Msg:    "block/unblock user success.",
		Error:  "",
	})
	c.Next()

}

func SetCommentVisibility(c *gin.Context) {
	targetComid := c.Query("comment_id")
	visibilityString := c.Query("visibility")

	visibility, err := strconv.Atoi(visibilityString)
	result := utils.DB.Where("comid = ?", targetComid).Model(models.Comment{}).Update("visibility", visibility)
	if err != nil || result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem setting comment visibility.",
			Error:  fmt.Sprintf("Having problem setting comment visibility, error: %v， %v", err, result.Error),
		})
		log.Printf("Having problem setting comment visibility, error: %v， %v", err, result.Error)
		c.Abort()
		return
	}
	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   "set comment visibility success",
		Msg:    "set comment visibility success",
		Error:  "",
	})
	c.Next()
}

func SetVideoVisibility(c *gin.Context) {
	targetVid := c.Query("vid")
	visibilityString := c.Query("visibility")

	visibility, err := strconv.Atoi(visibilityString)
	result := utils.DB.Where("vid = ?", targetVid).Model(models.Video{}).Update("visibility", visibility)
	if err != nil || result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem setting video visibility.",
			Error:  fmt.Sprintf("Having problem setting video visibility, error: %v， %v", err, result.Error),
		})
		log.Printf("Having problem setting video visibility, error: %v， %v", err, result.Error)
		c.Abort()
		return
	}
	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   "set video visibility success",
		Msg:    "set video visibility success",
		Error:  "",
	})
	c.Next()
}
