package controllerV1

import (
	"Server/app/constants"
	"Server/app/model"
	"Server/app/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"github.com/arangodb/go-driver"
)

//Create new tournament
func CreateTournament(c *gin.Context) {

	//Validation
	var tournament model.Tournament
	if err := c.ShouldBind(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	t := time.Now().Unix()

	ticket, _ := strconv.Atoi(c.PostForm("ticket"))
	totalCount, _ := strconv.Atoi(c.PostForm("total_count"))

	if ticket == 0 {
		tournament.Sum = 0
		tournament.Percentage = "%0"
		tournament.Income = 0
		tournament.Award = 0
	} else {
		percent, _ := strconv.Atoi(c.PostForm("percentage_of_dividends"))
		total := ticket * totalCount
		tournament.Percentage = "%" + strconv.Itoa(percent)
		tournament.Award = total - (((ticket * totalCount) * percent) / 100)
		tournament.Income = total - tournament.Award
		tournament.Sum = total
	}

	tournament.Approved = true
	tournament.State = constants.RegistrationState
	tournament.CreatedAt = t
	tournament.UpdatedAt = t

	//Create new tournament
	tournamentID, err := model.Create(&tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در ایجاد مسابقه ی جدید"}})
		return
	}

	gameID := driver.NewDocumentID(constants.Games, c.PostForm("game_key"))

	//Create tournament_to_game relation in games_edge
	if _, err := model.Create(&model.GamesEdge{From: tournamentID, To: gameID, Type: constants.TournamentToGame}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در ایجاد ارتباط tournament to game"}})
		return
	}

	//Create game_to_tournament relation in games_edge
	if _, err := model.Create(&model.GamesEdge{From: gameID, To: tournamentID, Type: constants.GameToTournament}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در ایجاد ارتباط game to tournament"}})
		return
	}

	ruleID := driver.NewDocumentID(constants.Rules, c.PostForm("rule_key"))

	//Create tournament_to_rule relation in games_edge
	if _, err := model.Create(&model.GamesEdge{From: tournamentID, To: ruleID, Type: constants.TournamentToRule}); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در ایجاد ارتباط tournament to rule"}})
		return
	}

	c.JSON(http.StatusOK, &response.Data{
		Data: &response.CreateDocument{
			Message: "مسابقه با موفقیت ایجاد شد",
			Key:     tournamentID.Key(),
		},
		State: true,
	})
}
