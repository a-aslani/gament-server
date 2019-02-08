package controllerApiV1

import (
	"Server/app/constants"
	"Server/app/model"
	"Server/app/response"
	"Server/app/utility"
	"github.com/arangodb/go-driver"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//Create new team
func CreateTeam(c *gin.Context) {

	//JWT Data
	claims := c.MustGet(constants.TokenClaims).(jwt.MapClaims)
	userKey := claims["key"].(string)

	teamImage, _ := c.FormFile("image")

	//Validation team post data
	var team model.Team
	if err := c.ShouldBind(&team); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{
			Data: &response.ValidationError{ Error: err.Error() },
		})
		return
	}

	if teamImage != nil {

		isUploaded, path, errMsg := utility.UploadImageWithResize(teamImage, constants.ImageTeamFolder, 80, 80)

		if !isUploaded {
			c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: errMsg}})
			return
		}

		team.Image = path

	} else {

		team.Image = constants.DefaultTeamImagePath
	}

	team.Admin = driver.NewDocumentID(constants.Users, userKey)

	t := time.Now().Unix()
	team.CreatedAt = t
	team.UpdatedAt = t

	//Create new Team
	teamID, err := model.Create(&team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در ایجاد تیم جدید"}})
		return
	}

	gameID := driver.NewDocumentID(constants.Games, team.GameKey)

	//Create relation from team to game
	if _, err := model.Create(&model.TeamsEdge{From: teamID, To: gameID, Type: constants.TeamToGame}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
		return
	}

	//Create relation from game to team
	if _, err := model.Create(&model.TeamsEdge{From: gameID, To: teamID, Type: constants.GameToTeam}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
		return
	}

	//Create relation from user to team
	if _, err := model.Create(&model.TeamsEdge{From: team.Admin, To: teamID, Type: constants.UserToTeam}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
		return
	}

	//Create relation from team to user
	if _, err := model.Create(&model.TeamsEdge{From: teamID, To: team.Admin, Type: constants.TeamToUser}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.CreateDocument{
			Message: "تیم جدید با موفقیت ایجاد شد",
			Key: teamID.Key(),
		},
		State: true,
	})
}