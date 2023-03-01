package userControllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go_round4/models"
	"go_round4/models/response"
	"go_round4/res"
	"go_round4/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

const NameLogin = 0
const MailLogin = 1

func Register(c *gin.Context) {
	newUser := models.User{}
	name := c.PostForm("user_name")
	password := c.PostForm("password")
	email := c.PostForm("email")
	slogan := c.PostForm("slogan")
	file, _ := c.FormFile("user_profile_image")

	//检查重名
	var searchOutcome []models.User
	result := utils.DB.Where("name = ?", name).First(&searchOutcome)
	if fmt.Sprintf("%v", result.Error) != "record not found" && result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Server have problems searching databases.",
			Error:  fmt.Sprintf("Databases search unavailable, err: %v", result.Error),
		})
		log.Printf("Databases search unavailable, err: %v", result.Error)
		c.Abort()
		return
	}

	if name == "" || password == "" || email == "" {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "Null user name, email or password.",
			Error:  "Illegal input",
		})
		log.Printf("Illegal name, email or password: %s, %s, %s \n", name, password, email)
		c.Abort()
		return
	}

	if result.Error == nil {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "The user name has been used.",
			Error:  "Used user name",
		})
		log.Printf("Used user name")
		c.Abort()
		return
	}

	//创建id
	userId, err := uuid.NewUUID()
	if err != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem generating user id.",
			Error:  fmt.Sprintf("Uuid generation failed error: %v", err),
		})
		log.Printf("Uuid generation failed error: %v", err)
		c.Abort()
		return
	}

	//加密密码
	encryptedPasswordBytes := encryptPassword(c, password)
	if encryptedPasswordBytes == nil {
		return
	}

	newUser.Uid = userId.String()
	//存储用户头像
	err = setUserHead(c, file, &newUser, 0)
	if err != nil {
		return
	}

	//生成用户数据结构
	newUser.EncryptedPassword = string(encryptedPasswordBytes)
	newUser.Name = name
	newUser.Email = email
	newUser.IsBlocked = 0
	newUser.Slogan = slogan

	err = utils.SaveUser(c, &newUser)
	if err != nil {
		return
	}

	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   fmt.Sprintf("userId: %s", newUser.Uid),
		Msg:    "Registration success!",
		Error:  "",
	})
	log.Printf("new user registration success! userId: %s, name: %s",
		newUser.Uid, newUser.Name)
}

func Login(loginType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result *gorm.DB
		var userInfo string
		email := c.Query("user_email")
		name := c.Query("user_name")
		password := c.Query("password") //搜索用户

		var searchOutcome []models.User
		if loginType == NameLogin {
			userInfo = name
			result = utils.DB.Where("name = ?", name).First(&searchOutcome)
		} else if loginType == MailLogin {
			userInfo = email
			result = utils.DB.Where("email = ?", email).First(&searchOutcome)
		}

		if fmt.Sprintf("%v", result.Error) == "record not found" {
			//无此用户：
			c.JSON(500, response.StringDataResponse{
				Status: 500,
				Data:   "",
				Msg:    "No such user.",
				Error:  fmt.Sprintf("No such user. userInfo: %s", userInfo),
			})
			log.Printf("No such user. userInfo: %s.", userInfo)
			c.Abort()
			return
		}

		if result.Error != nil {
			c.JSON(502, response.StringDataResponse{
				Status: 502,
				Data:   "",
				Msg:    "Having problem searching user in databases.",
				Error:  fmt.Sprintf("User search error: %v, userInfo: %s", result.Error, userInfo),
			})
			log.Printf("User search error: %v, userInfo: %s", result.Error, userInfo)
			c.Abort()
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(searchOutcome[0].EncryptedPassword), []byte(password))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			//错误密码：
			c.JSON(500, response.StringDataResponse{
				Status: 500,
				Data:   "",
				Msg:    "Wrong password",
				Error:  "Wrong password",
			})
			log.Printf("Wrong password: user: %v", userInfo)
			c.Abort()
			return
		}

		if err == nil {
			// 签发 token
			token, err := utils.Sign(res.USER_MODE, searchOutcome[0].Uid, searchOutcome[0].Name)
			if err != nil {
				utils.ShowUnknownTokenError(c, err, userInfo, 1)
				c.Abort()
				return
			}

			c.JSON(200, response.JsonDataResponse{
				Status: 200,
				Data:   gin.H{"token": token, "userId": searchOutcome[0].Uid},
				Msg:    "Login success",
				Error:  "",
			})
			log.Printf("Login success: user: %v", userInfo)
		}
	}
}

func PutSlogan(c *gin.Context) {
	uid := c.Query("user_id")
	slogan := c.Query("slogan")

	//search user
	searchOutcome := utils.SearchUser(c, uid, 1)
	if searchOutcome == nil {
		return
	}

	searchOutcome[0].Slogan = slogan

	//save info
	err := utils.SaveUser(c, &searchOutcome[0])
	if err != nil {
		return
	}

	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   fmt.Sprintf("userId: %s", searchOutcome[0].Uid),
		Msg:    "Slogan change success!",
		Error:  "",
	})
	log.Printf("Slogan change success success! userId: %s", uid)
}

func GetAvatar(c *gin.Context) {
	//简写了，因为 token 签发必须存在的 uid，在更改前会校验 token 中的 uid 信息，如果此处出错不会漏。
	//todo 缺少检测默认头像
	uid := c.Query("user_id")

	c.File(res.UserAvatarPath + uid + ".jpg")
}

func SearchUser(c *gin.Context) { // 模糊搜索
	content := c.Query("search_content")

	var resultList []models.User
	var result *gorm.DB

	result = utils.DB.Model(&models.User{}).Where("name LIKE ?", "%"+content+"%").Find(&resultList)

	if result.Error != nil && fmt.Sprintf("%v", result.Error) != "record not found" {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem searching items.",
			Error:  fmt.Sprintf("Search items failed, error: %v", result.Error),
		})
		log.Printf("Search items failed, error: %v", result.Error)
		c.Abort()
		return
	}
	c.JSON(200, response.UserListResponse{
		Status: 200,
		Data:   resultList,
		Msg:    "Search users success!",
		Error:  "",
	})
	log.Printf("Search user success!")
}

func PostUserHead(c *gin.Context) {
	file, err := c.FormFile("user_head_profile")
	uid := c.PostForm("user_id")

	searchOutcome := utils.SearchUser(c, uid, 1)
	if searchOutcome == nil {
		return
	}
	user := searchOutcome[0]
	err = setUserHead(c, file, &user, 1)
	if err == nil {
		c.JSON(200, response.StringDataResponse{
			Status: 200,
			Data:   fmt.Sprintf("userId: %s", user.Uid),
			Msg:    "change user head success!",
			Error:  "",
		})
		log.Printf("change user head success! userId: %s", user.Uid)
	}

	err = utils.SaveUser(c, &user)
	if err != nil {
		return
	}
}

func PutUserEmail(c *gin.Context) {
	uid := c.Query("user_id")
	mail := c.Query("user_mail")
	searchOutcome := utils.SearchUser(c, uid, 1)
	if searchOutcome == nil {
		return
	}

	searchOutcome[0].Email = mail

	err := utils.SaveUser(c, &searchOutcome[0])
	if err != nil {
		return
	}

	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   fmt.Sprintf("userId: %s", searchOutcome[0].Uid),
		Msg:    "email change success!",
		Error:  "",
	})
	log.Printf("email change success success! userId: %s", uid)
}

func SetPassword(c *gin.Context) {
	uid := c.Query("user_id")
	oldPas := c.Query("old_password")
	newPas := c.Query("new_password")

	searchOutcome := utils.SearchUser(c, uid, 1)
	if searchOutcome == nil {
		return
	}

	//对比数据库旧密码和输入的旧密码：
	err := bcrypt.CompareHashAndPassword([]byte(searchOutcome[0].EncryptedPassword), []byte(oldPas))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		//错误密码：
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "Wrong password",
			Error:  "Wrong password",
		})
		log.Printf("Wrong password: user: %v", searchOutcome[0].Uid)
		c.Abort()
		return
	}
	if err != nil {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "Error comparing password",
			Error:  "Error comparing password",
		})
		log.Printf("Error comparing password: user: %v, e: %v", searchOutcome[0].Uid, err)
		c.Abort()
		return
	} else {
		searchOutcome[0].EncryptedPassword = string(encryptPassword(c, newPas))
		err = utils.SaveUser(c, &searchOutcome[0])
		if err != nil {
			return
		}

		c.JSON(200, response.StringDataResponse{
			Status: 200,
			Data:   "",
			Msg:    "change password success",
			Error:  "",
		})
		log.Printf("change password success: user: %v", searchOutcome[0].Uid)
	}
}

func GetUserInfo(c *gin.Context) {
	uid := c.Query("user_id")

	var resultUser models.User
	result := utils.DB.Where("uid = ?", uid).First(&resultUser)
	if result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "error getting user info.",
			Error:  fmt.Sprintf("error getting user info. err: %v", result.Error),
		})
		log.Printf("error getting user info. err: %v", result.Error)
		c.Abort()
		return
	}

	c.JSON(200, response.UserResponse{
		Status: 200,
		Data:   resultUser,
		Msg:    "get user info success.",
		Error:  fmt.Sprintf("get user info success"),
	})
	c.Next()
}

func GetCollections(c *gin.Context) {
	uid := c.Query("user_id")

	var collectionList []models.Collection
	var collectionNameList []models.CollectionName
	result := utils.DB.Where("uid = ?", uid).Find(&collectionList)
	result2 := utils.DB.Where("uid = ?", uid).Find(&collectionNameList)
	if result.Error != nil || result2.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "unable to fetch collection data.",
			Error:  fmt.Sprintf("unable to fetch collection data. err: %v, %v", result.Error, result2.Error),
		})
		log.Printf("unable to fetch collection data. err: %v, %v", result.Error, result2.Error)
		c.Abort()
		return
	}

	collectionNameMap := make(map[string]string)
	for _, collectionName := range collectionNameList {
		collectionNameMap[collectionName.Colid] = collectionName.Name
	}
	collectionCount := len(collectionNameList)
	fullCollectionList := make([]models.FullCollection, collectionCount)
	for i := 0; i < len(collectionList); i++ {
		tempCollection := collectionList[i]
		fullCollectionList[i] = models.FullCollection{
			Colid: tempCollection.Colid,
			Uid:   tempCollection.Uid,
			Vid:   tempCollection.Vid,
			Name:  collectionNameMap[tempCollection.Colid],
		}
	}

	c.JSON(200, response.CollectionListResponse{
		Status: 200,
		Data:   fullCollectionList,
		Msg:    "get collection data success.",
		Error:  fmt.Sprintf("get collection data success."),
	})
	c.Next()
	//todo 不知道怎么生成多重列表然后返回，但是目前这个也能用就是
}

func GetSearchHistories(c *gin.Context) {
	uid := c.Query("user_id")

	var resultList []models.SearchHistory
	result := utils.DB.Where("uid = ?", uid).Order("search_time").Find(&resultList)

	if result.Error != nil && fmt.Sprintf("%v", result.Error) != "record not found" {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "unable to fetch search result data.",
			Error:  fmt.Sprintf("unable to fetch search result data.. err: %v", result.Error),
		})
		log.Printf("unable to fetch search result data.. err: %v", result.Error)
		c.Abort()
		return
	}

	c.JSON(200, response.SearchHistoryListResponse{
		Status: 200,
		Data:   resultList,
		Msg:    "get search history success.",
		Error:  "get search history success.",
	})
	c.Next()
}

// ------------ internal functions ------------

func encryptPassword(c *gin.Context, password string) []byte {
	encryptedPasswordBytes, err2 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err2 != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem generating user info.",
			Error:  fmt.Sprintf("EncryptedPswd generation failed, error: %v", err2),
		})
		log.Printf("EncryptedPswd generation failed, error: %v", err2)
		c.Abort()
		return nil
	}
	return encryptedPasswordBytes
}

func setUserHead(c *gin.Context, file *multipart.FileHeader, user *models.User, mode int) error { // mode: 0-create 1-save
	if file != nil {
		//限制 JPG
		name := strings.Split(file.Filename, ".")
		if name[len(name)-1] != "jpg" {
			log.Printf("IP %s: file upload failed, error: %v", c.ClientIP(), errors.New(fmt.Sprintf("unsupported file format! use jpg instead")))
			c.JSON(500, response.StringDataResponse{
				Status: 500,
				Data:   "",
				Msg:    "header image upload failed - file type not supported",
				Error:  fmt.Sprintf("unsupported file format!"),
			})
			c.Abort()
			return errors.New("setUserHead error - filetype")
		}

		if user.Uid == "" {
			return errors.New("setUserHead error - empty user data")
		}

		dst := res.UserAvatarPath + user.Uid + ".jpg" //懒了，只允许jpg
		err := c.SaveUploadedFile(file, dst)
		//上传失败
		if err != nil {
			log.Printf("IP %s: file upload failed, error: %e", c.ClientIP(), err)
			c.JSON(502, response.StringDataResponse{
				Status: http.StatusBadGateway,
				Data:   "",
				Msg:    "header image upload failed - unknown error",
				Error:  fmt.Sprintf("%e", err),
			})
			c.Abort()
			return errors.New("setUserHead error - upload")
		}
		user.ProfileImageLocation = dst
	} else if mode == 0 {
		user.ProfileImageLocation = res.DefaultAvatarPath
	}

	return nil
}
