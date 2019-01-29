package controllerV1

import (
	"Server/app/constants"
	"Server/app/model"
	"Server/app/response"
	"Server/app/utility"
	"Server/src/src/github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Create new Game
func CreateGame(c *gin.Context) {

	imageFile, _ := c.FormFile("image")
	mobileBannerFile, _ := c.FormFile("mobile_banner")
	webBannerFile, _ := c.FormFile("web_banner")

	//Validation image
	if imageFile == nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "انتخاب تصویر ضروری است"}})
		return
	} else if mobileBannerFile == nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "انتخاب تصویر بنر نسخه ی موبایل ضروری است"}})
		return
	} else if webBannerFile == nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "انتخاب تصویر بنر نسخه ی وب ضروری است"}})
		return
	}

	//Validation
	var game model.Game
	if err := c.ShouldBind(&game); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	game.PC, _ = strconv.ParseBool(c.PostForm("pc"))
	game.PS, _ = strconv.ParseBool(c.PostForm("ps"))
	game.Xbox, _ = strconv.ParseBool(c.PostForm("xbox"))
	game.Mobile, _ = strconv.ParseBool(c.PostForm("mobile"))

	isUploadGameImage, gameImagePath, gameImageErrMsg := utility.UploadImageCustom(imageFile, constants.ImageGame, 350, 195)
	isUploadGameMobileImage, gameMobileImagePath, gameMobileImageErrMsg := utility.UploadImageCustom(mobileBannerFile, constants.ImageBanner, 350, 340)
	isUploadGameWebImage, gameWebImagePath, gameWebImageErrMsg := utility.UploadImageCustom(webBannerFile, constants.ImageBanner, 1600, 240)


	if !isUploadGameImage {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: gameImageErrMsg}})
		return
	} else if !isUploadGameMobileImage {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: gameMobileImageErrMsg}})
		return
	} else if !isUploadGameWebImage {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: gameWebImageErrMsg}})
		return
	}


	t := time.Now().Unix()
	game.Image = gameImagePath
	game.MobileBanner = gameMobileImagePath
	game.WebBanner = gameWebImagePath
	game.Approved = true
	game.CreatedAt = t
	//game.UpdatedAt = t

	//Create new game
	gameID, err := model.Create(&game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در ایجاد بازی جدید"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.CreateDocument{
			Message: "بازی با موفقیت ایجاد شد",
			Key:     gameID.Key(),
		},
		State: true,
	})
}

//Find game
func FindGame(c *gin.Context) {

	key := c.Param("key")

	game, err := model.Find(key, &model.Game{})

	if driver.IsNotFound(err) || err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "بازی پیدا نشد"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.FindDocument{
			Document: game,
		},
		State: true,
	})
}

//Find all games
func FindAllGames(c *gin.Context) {

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 0, 64)

	games, err := model.FindAll(constants.GamesCount, page, &model.Game{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.DatabaseError{Message: "بازی ایی در دیتابیس یافت نشد"}})
		return
	}

	if games == nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "بازی ایی وجود ندارد"}})
		return
	}

	count, err := model.Count(&model.Game{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.DatabaseError{Message: "بازی پیدا نشد"}})
		return
	}

	pages := utility.Pages(count, constants.GamesCount)

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.FindAllDocuments{
			Documents:  games,
			TotalPages: pages,
			CurrentPage: page,
		},
		State: true,
	})
}

//Update game
func UpdateGame(c *gin.Context) {

	imageFile, _ := c.FormFile("image")
	mobileBannerFile, _ := c.FormFile("mobile_banner")
	webBannerFile, _ := c.FormFile("web_banner")
	key := c.Param("key")

	var game model.Game
	if err := c.ShouldBind(&game); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	gameDoc, err := model.Find(key, &model.Game{})
	if driver.IsNotFound(err) || err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "بازی پیدا نشد"}})
		return
	}

	if imageFile != nil {
		isUpload, path, errMsg := utility.UploadImage(imageFile, constants.ImageGame)
		if !isUpload {
			c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: errMsg}})
			return
		}
		//Upload image
		c.SaveUploadedFile(imageFile, path)
		os.Remove(gameDoc["image"].(string))
		game.Image = path
	} else {
		game.Image = gameDoc["image"].(string)
	}

	if mobileBannerFile != nil {
		isUpload, path, errMsg := utility.UploadImage(mobileBannerFile, constants.ImageBanner)
		if !isUpload {
			c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: errMsg}})
			return
		}
		//Upload image
		c.SaveUploadedFile(mobileBannerFile, path)
		os.Remove(gameDoc["mobile_banner"].(string))
		game.MobileBanner = path
	} else {
		game.MobileBanner = gameDoc["mobile_banner"].(string)
	}

	if webBannerFile != nil {
		isUpload, path, errMsg := utility.UploadImage(webBannerFile, constants.ImageBanner)
		if !isUpload {
			c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: errMsg}})
			return
		}
		//Upload image
		c.SaveUploadedFile(webBannerFile, path)
		os.Remove(gameDoc["web_banner"].(string))
		game.WebBanner = path
	} else {
		game.WebBanner = gameDoc["web_banner"].(string)
	}

	if c.PostForm("active") != "" {
		game.Approved, _ = strconv.ParseBool(c.PostForm("approved"))
	}
	if c.PostForm("pc") != "" {
		game.PC, _ = strconv.ParseBool(c.PostForm("pc"))
	}
	if c.PostForm("ps") != "" {
		game.PS, _ = strconv.ParseBool(c.PostForm("ps"))
	}
	if c.PostForm("mobile") != "" {
		game.Mobile, _ = strconv.ParseBool(c.PostForm("mobile"))
	}

	t := time.Now().Unix()
	game.CreatedAt = int64(gameDoc["created_at"].(float64))
	game.UpdatedAt = t

	//Update game
	err = model.Update(key, &game)
	if driver.IsNotFound(err) || err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "بازی پیدا نشد. خطا در ویرایش بازی"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.UpdateDocument{
			Message: "بازی با موفقیت ویرایش شد",
			Key:     key,
		},
		State: true,
	})
}

//Remove game
func DestroyGame(c *gin.Context) {
	key := c.Param("key")
	err := model.Destroy(key, &model.Game{})
	if driver.IsNotFound(err) || err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "بازی پیدا نشد. خطا در حذف کردن بازی"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.DestroyDocument{
			Message: "بازی با موفقیت حذف شد",
		},
		State: true,
	})

}
