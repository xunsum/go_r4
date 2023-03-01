package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_round4/models"
	"go_round4/models/response"
	"gorm.io/gorm"
	"log"
)

func ShowUnknownTokenError(c *gin.Context, err error, userName string, stampPoint int) {
	c.JSON(502, response.StringDataResponse{
		Status: 502,
		Data:   "",
		Msg:    "Having problem giving token.",
		Error:  fmt.Sprintf("Having problem giving token. error: %v, userName: %s", err, userName),
	})
	log.Printf("Having problem giving token. error: %v, userName: %s, stampPoint: %d", err, userName, stampPoint)
}

func SearchUser(c *gin.Context, searchContent string, mode int) []models.User { // mode: name - 0 id - 1
	var searchOutcome []models.User
	var result *gorm.DB
	if mode == 0 {
		result = DB.Where("uid = ?", searchContent).First(&searchOutcome)
	} else if mode == 1 {
		result = DB.Where("uid = ?", searchContent).First(&searchOutcome)
	}
	if result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Server have problems searching databases.",
			Error:  fmt.Sprintf("Databases search unavailable, err: %v", result.Error),
		})
		log.Printf("Databases search unavailable, err: %v", result.Error)
		c.Abort()
		return nil
	}
	return searchOutcome
}

func SaveUser(c *gin.Context, user *models.User) error {
	//存储用户信息
	result := DB.Where("uid = ?", user.Uid).Save(&user)
	if result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem saving user info.",
			Error:  fmt.Sprintf("Create user info failed, error: %v", result.Error),
		})
		log.Printf("Create user info failed, error: %v", result.Error)
		c.Abort()
		return errors.New("having problem saving user info")
	}
	return nil
}

// SaveVideoInfo
// 新建/更改视频信息
func SaveVideoInfo(c *gin.Context, video *models.Video) error {
	//添加/保存视频数据
	result := DB.Where("vid = ?", video.Vid).Save(&video)
	if result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem saving video info.",
			Error:  fmt.Sprintf("Create video info failed, error: %v", result.Error),
		})
		log.Printf("Create video info failed, error: %v", result.Error)
		c.Abort()
		return errors.New("having problem saving video info")
	}
	return nil
}

// SetLike
// vid:Video that's watched,
// uid: The user who's watching,
// isLike:
func SetLike(c *gin.Context, vid string, uid string, isLike bool) error {

	//流程：搜索视频，点赞列表搜索点赞 -> 校验是否能够添加/减少点赞，搜索视频数据，增加点赞计数，添加/删除点赞项
	var resultVideo models.Video
	var resultLikeInfo models.Like

	result1 := DB.Model(&models.Video{}).Where("vid = ?", vid).First(&resultVideo)
	result2 := DB.Model(&models.Like{}).Where("vid = ?", vid).Where("uid = ?", uid).First(&resultLikeInfo)

	// Err when searching video DB
	if result1.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem fetching video.",
			Error:  fmt.Sprintf("Fetching video info failed, error: %v", result1.Error),
		})
		log.Printf("Fetching video info failed, error: %v", result1.Error)
		c.Abort()
		return errors.New("having problem fetching video info")
	}

	// Err when searching likes DB -> It's legal to have one when user didn't like the video.
	if (fmt.Sprintf("%v", result2.Error) != "record not found" && result2.Error != nil) ||
		(fmt.Sprintf("%v", result2.Error) == "record not found" && !isLike) ||
		(result2.Error == nil && isLike) {
		//When err is not no record, or we're unliking the video while there's no like,
		//or we're liking while there's already a like, real error happened:
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "Like/dislike info doesn't match DB info or like info search internal error",
			Error:  fmt.Sprintf("Like/dislike info doesn't match DB info or like info search internal error, error: %v", result2.Error),
		})
		log.Printf("Like/dislike info doesn't match DB info or like info search internal error, error: %v", result2.Error)
		c.Abort()
		return errors.New("like/dislike info doesn't match DB info or like info search internal error")
	}

	isSuccess := false
	var errCarrier error
	//liking:
	if fmt.Sprintf("%v", result2.Error) == "record not found" && isLike {
		//添加like记录，更改视频记录
		result := DB.Where("uid = ? ", uid).Where("vid = ?", vid).Save(&models.Like{Uid: uid, Vid: vid})
		resultVideo.LikeCount = resultVideo.LikeCount + 1 // todo: 肯定有更好的解决方案，理论上是耗时操作，同时也不能支持多线程
		err := SaveVideoInfo(c, &resultVideo)
		if err != nil {
			return errors.New("saving video info err")
		}
		if result.Error == nil {
			isSuccess = true
		} else {
			errCarrier = result.Error
		}
	}

	//unliking
	if result2.Error == nil && !isLike {
		//删除，更改视频记录
		result := DB.Where("uid = ? ", uid).Where("vid = ?", vid).Where("uid = ?", uid).Delete(&resultLikeInfo)
		resultVideo.LikeCount = resultVideo.LikeCount - 1 // todo: 肯定有更好的解决方案，理论上是耗时操作，同时也不能支持多线程
		err := SaveVideoInfo(c, &resultVideo)
		if err != nil {
			return errors.New("saving video info err")
		}
		if result.Error == nil {
			isSuccess = true
		} else {
			errCarrier = result.Error
		}
	}
	if !isSuccess {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Like/dislike error.",
			Error:  fmt.Sprintf("Like/dislike error. Err: %v", errCarrier),
		}) // todo: 这里最好返一个video结构体，方便前端更新，懒了
		log.Printf("Like/dislike error. Err: %v", errCarrier)
		c.Abort()
		return errors.New("like/dislike process error")
	}

	return nil
}

func SaveComment(c *gin.Context, comment models.Comment) {

}
