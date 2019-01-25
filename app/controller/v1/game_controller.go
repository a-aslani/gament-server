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

	//Validation image
	if imageFile == nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: "انتخاب تصویر ضروری است"}})
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

	isUpload, path, errMsg := utility.UploadImageCustom(imageFile, constants.ImageGame, 350, 195)

	if !isUpload {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: errMsg}})
		return
	}


	t := time.Now().Unix()
	game.Image = path
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
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "Game not found"}})
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
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.DatabaseError{Message: "cannot find games from database"}})
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
	key := c.Param("key")

	var game model.Game
	if err := c.ShouldBind(&game); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	gameDoc, err := model.Find(key, &model.Game{})
	if driver.IsNotFound(err) || err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "game not found, find error"}})
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
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "game not found, update error"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.UpdateDocument{
			Message: "Game successfully updated",
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
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "game not found, remove error"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.DestroyDocument{
			Message: "Game successfully removed",
		},
		State: true,
	})

}
