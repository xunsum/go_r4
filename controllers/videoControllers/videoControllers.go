package videoControllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go_round4/models"
	"go_round4/models/response"
	"go_round4/res"
	"go_round4/utils"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UploadVideo(c *gin.Context) {
	uid := c.PostForm("user_id")
	title := c.PostForm("title")
	description := c.PostForm("description")
	videoType := c.PostForm("type")
	length := c.PostForm("vid_length")

	//获取视频头
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "having problem fetching file!",
			Error:  "having problem fetching file!",
		})
		c.Abort()
		return
	}

	//创建数据
	vid, err := uuid.NewUUID()
	if err != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem generating vid.",
			Error:  fmt.Sprintf("Uuid generation failed error: %v", err),
		})
		log.Printf("Uuid generation failed error: %v", err)
		c.Abort()
		return
	}

	uploadTime := strconv.FormatInt(time.Now().Unix(), 10)

	newVideo := models.Video{
		LikeCount:   0,
		Vid:         vid.String(),
		Uid:         uid,
		Title:       title,
		Description: description,
		Type:        videoType,
		UploadTime:  uploadTime,
		Length:      length,
		Visibility:  1,
		Views:       0,
	}

	err = saveVideo(c, file, &newVideo)
	if err != nil {
		return
	}

	err = utils.SaveVideoInfo(c, &newVideo)
	if err != nil {
		return
	}

	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   fmt.Sprintf("vid: %s", vid.String()),
		Msg:    "Video upload success!",
		Error:  "",
	})
	log.Printf("Video upload success! vid: %s", vid)

}

func SetLikeVideo(c *gin.Context) {
	isLike := c.Query("like") == "1"
	uid := c.Query("user_id")
	vid := c.Query("vid")
	err := utils.SetLike(c, vid, uid, isLike)
	if err != nil {
		return
	}

	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   "",
		Msg:    "Like/dislike success.",
		Error:  fmt.Sprintf("Like/dislike success."),
	})
}

func SetComment(c *gin.Context) {
	uid := c.Query("user_id")
	vid := c.Query("vid")
	content := c.Query("content")
	replyTo := c.Query("reply_to")

	comid, err := uuid.NewUUID()
	if err != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem generating vid.",
			Error:  fmt.Sprintf("Uuid generation failed error: %v", err),
		})
		log.Printf("Uuid generation failed error: %v", err)
		c.Abort()
		return
	}

	newComment := models.Comment{
		Comid:      comid.String(),
		Uid:        uid,
		Vid:        vid,
		Content:    content,
		Visibility: "1",
		ReplyTo:    replyTo,
	}

	if replyTo != "" {
		//检测 reply_to 是否存在：
		targetComment := models.Comment{}
		result := utils.DB.Where("comid = ?", replyTo).First(&targetComment)
		if fmt.Sprintf("%v", result.Error) == "record not found" || targetComment.Vid != newComment.Vid {
			c.JSON(502, response.StringDataResponse{
				Status: 502,
				Data:   "",
				Msg:    "Having problem replying to comment.",
				Error:  fmt.Sprintf("Unable to fetch the comment replying to, not under the same video or doesn't exist."),
			})
			log.Printf("Unable to fetch the comment replying to, not under the same video or doesn't exist.")
			c.Abort()
			return
		}
	}

	result := utils.DB.Where("comid = ?", comid).Save(&newComment)
	if result.Error != nil {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    "Having problem creating comment.",
			Error:  fmt.Sprintf("Having problem creating comment: %v", result.Error),
		})
		log.Printf("Having problem creating comment: %v", result.Error)
		c.Abort()
		return
	}

	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   "",
		Msg:    "Set comment success!",
		Error:  fmt.Sprintf("Set comment success!"),
	})
	c.Next()
}

func SetCollection(c *gin.Context) {
	//情况:
	//1. 新建空：没有 collectionId，没有 vid，直接新建 collectionId，有 vid 则首个视频是此 vid，否则生成占位的 vid = 0
	//2. 增加：有 collectionId，同时有 vid，直接save就行了
	//3. 新建同时增加：没有 collectionId，有 vid，创建Id然后 save
	//4. 删除收藏内容：operation = 0
	//5. 删除收藏夹：op = 0, vid = ""
	//6. 更改收藏夹名字：op = "3"，必须有 collection 信息，只看 collection_name

	uid := c.Query("user_id")
	vid := c.Query("vid")
	colid := c.Query("collection_id")
	name := c.Query("collection_name")
	operation := c.Query("operation")

	var result1 *gorm.DB //对收藏数据进行操作的结果
	var result2 *gorm.DB //对收藏名称数据库操作的结果
	var optionHint string

	//新建
	if operation == "1" {

		//禁止自定义collectionId，因为使用save指令简写，所以可能通过这种方式自定义ID：
		result1 = utils.DB.Where("colid = ?", colid).First(&models.Collection{})
		if fmt.Sprintf("%v", result1.Error) == "record not found" {
			colid = ""
		}

		//没有cid/vid就创建:
		if colid == "" {
			//没有 colid，肯定是创建新的收藏夹（1/3）
			tempColid, err := uuid.NewUUID()
			if err != nil {
				c.JSON(502, response.StringDataResponse{
					Status: 502,
					Data:   "",
					Msg:    "Having problem generating vid.",
					Error:  fmt.Sprintf("Uuid generation failed error: %v", err),
				})
				log.Printf("Uuid generation failed error: %v", err)
				c.Abort()
				return
			}
			//没有 vid 就创建占位符（1）
			if vid == "" {
				vid = "1"
			}
			colid = tempColid.String()
			optionHint = "create collection"
		} else if vid != "" {
			//有 colid，同时有 vid，是添加新的视频（2）
			optionHint = "add vid"
		}

		//新建/添加 - 数据库操作：
		newCollectionItem := models.Collection{Uid: uid, Vid: vid, Colid: colid}
		newCollectionName := models.CollectionName{Name: name, Colid: colid, Uid: uid}
		result1 = utils.DB.Where("colid = ?", colid).Where("vid = ?", vid).Save(&newCollectionItem)
		result2 = utils.DB.Where("colid = ?", colid).Where("name = ?", name).Save(&newCollectionName)

		optionHint = "add/creating collection/video"
	}

	//更改名
	if operation == "2" {
		//检测是否存在：
		result2 = utils.DB.Where("colid = ?", colid).First(&models.CollectionName{})
		var err error
		if fmt.Sprintf("%v", result2.Error) == "record not found" {
			err = errors.New("collection/collectionName not found")
		}
		//名字不能为空
		if name == "" {
			err = errors.New("name of changing collection is not provided")
		}

		//条件过不过：
		if err != nil {
			c.JSON(500, response.StringDataResponse{
				Status: 500,
				Data:   "",
				Msg:    "Unable to change collection name.",
				Error:  fmt.Sprintf("Unable to change collection name: %v", err),
			})
			log.Printf("Unable to change collection name: %v", err)
			c.Abort()
			return
		}

		//更改
		result2 = utils.DB.Where("colid = ?", colid).Save(&models.CollectionName{
			Colid: colid,
			Name:  name,
			Uid:   uid,
		})
		optionHint = "changing collection name"
		result1 = result2
	}

	//删除
	if operation == "0" {
		if vid != "" {
			var tempCollist []models.Collection
			//删除后如果什么都不剩，应该添加占位符
			result1 = utils.DB.Where("colid = ?", colid).Where("vid = ?", vid).Find(&tempCollist)
			//删除一条视频
			if len(tempCollist) == 1 {
				//直接把最后一条改成占位符
				result1 = utils.DB.Where("colid = ?", colid).Where("vid = ?", vid).Save(&models.Collection{Colid: colid, Uid: uid, Vid: "1"})
			} else {
				//还有剩的
				result1 = utils.DB.Where("colid = ?", colid).Where("vid = ?", vid).Delete(models.Collection{})
			}
		} else {
			//删除所有的 collection
			result1 = utils.DB.Where("colid = ?", colid).Delete(models.Collection{})
			result2 = utils.DB.Where("colid = ?", colid).Delete(models.CollectionName{})
		}
		optionHint = "deleting collection/video"
	}

	//其他数字，直接错误
	if operation != "1" && operation != "0" && operation != "2" {
		optionHint = "error"
	}

	//错误处理
	if result1.Error != nil || result2.Error != nil || optionHint == "error" {
		c.JSON(502, response.StringDataResponse{
			Status: 502,
			Data:   "",
			Msg:    fmt.Sprintf("Having problem %v", optionHint),
			Error:  fmt.Sprintf("Having problem %v, err1: %v, err2: %v", optionHint, result1.Error, result2.Error),
		})
		log.Printf("Having problem %v, err1: %v, err2: %v", optionHint, result1.Error, result2.Error)
		c.Abort()
		return
	}

	c.JSON(200, response.StringDataResponse{
		Status: 200,
		Data:   "",
		Msg:    "collection opt success!",
		Error:  fmt.Sprintf("collection opt： %v success!", optionHint),
	})
	c.Next()
	return
}

func SearchVideo(c *gin.Context) {
	uid := c.Query("user_id")
	userName := c.Query("user_name")
	title := c.Query("title")
	description := c.Query("description")
	videoType := c.Query("type")
	uploadTime := c.Query("uploadTime")  //范围
	likeCounts := c.Query("like_counts") //范围
	views := c.Query("view_counts")      //范围
	length := c.Query("length")          //范围

	resultList := []models.Video{}
	var errList []error

	//by name:
	if userName != "" {
		var userList []models.User
		var subList []models.Video
		//get uids
		db1 := utils.DB.Where("name LIKE ?", "%"+userName+"%")
		result1 := db1.Find(&userList)
		if result1.Error != nil {
			errList = append(errList, errors.New("name search err1"))
		}
		//get all video infos
		for i := 0; i < len(userList); i++ {
			var tempList []models.Video
			db2 := utils.DB.Where("uid = ?", userList[i].Uid)
			processLimits(&errList, db2, videoType, uploadTime, likeCounts, views, length)
			result2 := db2.Find(&tempList)
			subList = append(subList, tempList...)
			if result2.Error != nil {
				errList = append(errList, errors.New("name search err2"))
				break
			}
		}
		//append to main list
		resultList = append(resultList, subList...)
	}

	//by title:
	if title != "" {
		var subList []models.Video
		db := utils.DB.Where("title LIKE ?", "%"+title+"%")
		processLimits(&errList, db, videoType, uploadTime, likeCounts, views, length)
		result := db.Find(&subList)
		if result.Error != nil {
			errList = append(errList, errors.New("title search err"))
		}
		resultList = append(resultList, subList...)
	}

	//by description
	if description != "" {
		var subList []models.Video
		db := utils.DB.Where("description LIKE ?", "%"+description+"%")
		processLimits(&errList, db, videoType, uploadTime, likeCounts, views, length)
		result := db.Find(&subList)
		if result.Error != nil {
			errList = append(errList, errors.New("description search err"))
		}
		resultList = append(resultList, subList...)
	}

	if len(errList) != 0 && len(resultList) == 0 {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "err occurred while searching for video, no result was found",
			Error:  fmt.Sprintf("err occurred while searching for video, no result was found, errs: %v", errList),
		})
		log.Printf("err occurred while searching for video, no result was found, errs: %v", errList)
		c.Abort()
		return
	} else if len(errList) != 0 && len(resultList) != 0 {
		c.JSON(200, response.VideoListResponse{
			Status: 200,
			Data:   resultList,
			Msg:    "some err occurred, but search is half success",
			Error:  fmt.Sprintf("some err occurred, but search is half success, err: %v", errList),
		})
		log.Printf("some err occurred, but search is half success, err: %v", errList)
		c.Next()
	} else if len(errList) == 0 {
		c.JSON(200, response.VideoListResponse{
			Status: 200,
			Data:   resultList,
			Msg:    "search success",
			Error:  fmt.Sprintf("search success"),
		})
		c.Next()
	}

	//保存历史
	searchTime := time.Now().Unix()
	var searchType int
	tempType, err := strconv.Atoi(videoType)
	if err == nil {
		searchType = tempType
	}
	searchHistory := models.SearchHistory{
		SearchTime:  searchTime,
		Uid:         uid,
		Name:        userName,
		Title:       title,
		Description: description,
		Type:        searchType,
		UploadTime:  uploadTime,
		LikeCounts:  likeCounts,
		ViewCounts:  views,
		Length:      length,
	}
	result := utils.DB.Where("search_time = ?", searchTime).Where("uid = ?", uid).Save(&searchHistory)
	if result.Error != nil {
		log.Printf("save search history error, err: %v", result.Error)
	}
}

func SetDanmaku(c *gin.Context) {
	uid := c.Query("user_id")
	dmkContent := c.Query("danmaku_content")
	vid := c.Query("video_id")
	danmakuTime, err := strconv.Atoi(c.Query("danmaku_time"))

	for {
		if err != nil {
			break
		}
		//验证视频是否存在：
		result := utils.DB.Where("vid = ?", vid).First(&models.Video{})
		if fmt.Sprintf("%v", result.Error) == "record not found" {
			err = errors.New("set danmaku video not found")
			break
		}

		newDid, tmepErr := uuid.NewUUID()
		if tmepErr != nil {
			err = errors.New("set danmaku: did generation failed")
			break
		}

		newDanmaku := models.Danmaku{
			Uid:      uid,
			Did:      newDid.String(),
			Vid:      vid,
			Content:  dmkContent,
			SentTime: danmakuTime,
		}

		result = utils.DB.Where("did = ?", newDid).Save(&newDanmaku)
		fmt.Printf("------------------------------> %v \n", result.Error)
		if result.Error != nil {
			err = errors.New("set danmaku: save danmaku in db failed")
		}

		break
	}

	if err != nil {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "err occurred while setting danmaku",
			Error:  fmt.Sprintf("err occurred while setting danmaku, err: %v", err),
		})
		log.Printf("err occurred while setting danmaku, err: %v", err)
		c.Abort()
		return
	} else {
		c.JSON(200, response.StringDataResponse{
			Status: 200,
			Data:   "",
			Msg:    "set danmaku success",
			Error:  fmt.Sprintf("search success"),
		})
		c.Next()
		return
	}
}

func GetVideo(c *gin.Context) {
	vid := c.Query("vid")
	//验证存在：
	result := utils.DB.Where("vid = ?").First(models.Video{})
	if fmt.Sprintf("%v", result.Error) == "record not found" {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "Failed to find the video",
			Msg:    "Failed to find the video",
			Error:  fmt.Sprintf("Failed to find the video, err: %v", result.Error),
		})
		log.Printf("Failed to find the video, err: %v", result.Error)
		c.Abort()
		return
	}
	c.File(res.VideoPath + vid + ".mp4")

}

func GetDanmakuList(c *gin.Context) {
	vid := c.Query("vid")

	var damakuList []models.Danmaku
	result := utils.DB.Order("sent_time").Where("vid = ?", vid).Find(&damakuList)
	if result.Error != nil {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "err occurred while getting danmakuList",
			Error:  fmt.Sprintf("err occurred while getting danmakuList, err: %v", result.Error),
		})
		log.Printf("err occurred while getting danmakuList, err: %v", result.Error)
		c.Abort()
		return
	}

	c.JSON(200, response.DanmakuListResponse{
		Status: 200,
		Data:   damakuList,
		Msg:    "get danmakuList success",
		Error:  fmt.Sprintf("get danmakuList success"),
	})
	c.Next()
	return
}

func GetVideoInfo(c *gin.Context) {
	vid := c.Query("vid")

	var videoInfo models.Video
	result := utils.DB.Where("vid = ?", vid).First(&videoInfo)

	if result.Error != nil {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "err occurred while getting video info",
			Error:  fmt.Sprintf("err occurred while getting video info, err: %v", result.Error),
		})
		log.Printf("err occurred while getting video info, err: %v", result.Error)
		c.Abort()
		return
	} else {
		c.JSON(200, response.VideoResponse{
			Status: 200,
			Data:   videoInfo,
			Msg:    "get video info success.",
			Error:  "",
		})
		c.Next()
	}
}

func GetComments(c *gin.Context) {
	vid := c.Query("vid")

	var commentList []models.Comment
	result := utils.DB.Where("vid = ?", vid).Find(&commentList)
	if result.Error != nil {
		c.JSON(500, response.StringDataResponse{
			Status: 500,
			Data:   "",
			Msg:    "err occurred while getting commentList",
			Error:  fmt.Sprintf("err occurred while getting commentList, err: %v", result.Error),
		})
		log.Printf("err occurred while getting commentList, err: %v", result.Error)
		c.Abort()
		return
	}

	c.JSON(200, response.CommentListResponse{
		Status: 200,
		Data:   commentList,
		Msg:    "get comment success",
		Error:  fmt.Sprintf("get comment success"),
	})
	c.Next()
}

// ------------ internal functions ------------
// 用作上传参考
func saveVideo(c *gin.Context, file *multipart.FileHeader, video *models.Video) error { // mode: 0-create 1-save
	if file != nil {
		//限制 mp4
		name := strings.Split(file.Filename, ".")
		if name[len(name)-1] != "mp4" {
			log.Printf("IP %s: file upload failed, error: %v", c.ClientIP(), errors.New(fmt.Sprintf("unsupported file format! use mp4 instead")))
			c.JSON(500, response.StringDataResponse{
				Status: 500,
				Data:   "",
				Msg:    "video upload failed - file type not supported",
				Error:  fmt.Sprintf("unsupported file format! Use mp4 instead."),
			})
			c.Abort()
			return errors.New("setUserHead error - filetype")
		}

		if video.Vid == "" {
			return errors.New("setUserHead error - empty vid")
		}

		dst := res.VideoPath + video.Vid + ".mp4"
		err := c.SaveUploadedFile(file, dst)
		//上传失败
		if err != nil {
			log.Printf("IP %s: file upload failed, error: %e", c.ClientIP(), err)
			c.JSON(502, response.StringDataResponse{
				Status: http.StatusBadGateway,
				Data:   "",
				Msg:    "video upload failed - unknown error",
				Error:  fmt.Sprintf("%e", err),
			})
			c.Abort()
			return errors.New("upload Video error - upload")
		}
	}

	return nil
}
func resolveLimit(limitString string) (int, int, error) {
	var err, err1, err2 error
	var lowerLimit int
	var upperLimit int
	for true {
		limitArray := strings.Split(limitString, "-")
		if len(limitArray) != 2 {
			err = errors.New("unable to resolve limit info: split")
			break
		}

		lowerLimit, err1 = strconv.Atoi(limitArray[0])
		upperLimit, err2 = strconv.Atoi(limitArray[1])
		if err1 != nil || err2 != nil {
			err = errors.New("unable to resolve limit info: Atoi")
			break
		}

		break
	}

	return lowerLimit, upperLimit, err
}
func processLimits(errList *[]error, db *gorm.DB, videoType string, uploadTime string, likeCounts string, views string, length string) {
	if videoType != "" {
		db = db.Where("type = ?", videoType)
	}

	if uploadTime != "" {
		for {
			lowerBound, upperBound, err := resolveLimit(uploadTime)
			if err != nil {
				*errList = append(*errList, errors.New("uploadTime search error: limit conversion"))
				break
			}
			db = db.Where("upload_time > ?", lowerBound).Where("upload_time < ?", upperBound)
			break
		}
	}

	if likeCounts != "" {
		for {
			lowerBound, upperBound, err := resolveLimit(likeCounts)
			if err != nil {
				*errList = append(*errList, errors.New("likeCounts search error: limit conversion"))
				break
			}
			db = db.Where("like_count > ?", lowerBound).Where("like_count < ?", upperBound)
			break
		}
	}

	if views != "" {
		for {
			lowerBound, upperBound, err := resolveLimit(views)
			if err != nil {
				*errList = append(*errList, errors.New("views search error: limit conversion"))
				break
			}
			db = db.Where("views > ?", lowerBound).Where("views < ?", upperBound)
			break
		}
	}

	if length != "" {
		for {
			lowerBound, upperBound, err := resolveLimit(length)
			if err != nil {
				*errList = append(*errList, errors.New("length search error: limit conversion"))
				break
			}
			db = db.Where("length > ?", lowerBound).Where("length < ?", upperBound)
			break
		}
	}
}
