package controllerV1

import (
	"Server/app/constants"
	"Server/app/model"
	"Server/app/model/database"
	"Server/app/response"
	"Server/app/utility"
	"Server/app/utility/jalali"
	"github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Create new tournament
func CreateTournament(c *gin.Context) {

	ticketForm := c.PostForm("ticket")

	//Validation
	var tournament model.Tournament

	if ticketForm == "" {
		tournament.Ticket = 0
	} else {
		tournament.Ticket, _ = strconv.Atoi(ticketForm)
	}

	if err := c.ShouldBind(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, &response.Data{Data: &response.ValidationError{Error: err.Error()}})
		return
	}

	//Platform to upper case
	tournament.Platform = strings.ToUpper(strings.TrimSpace(c.PostForm("platform")))

	t := time.Now().Unix()

	ticket := tournament.Ticket
	quantity := tournament.Quantity

	if ticket == 0 {
		tournament.Sum = 0
		tournament.Percentage = "%0"
		tournament.Income = 0
		tournament.Award = 0
	} else {
		percent, _ := strconv.Atoi(c.PostForm("percentage"))
		total := ticket * quantity
		tournament.Percentage = "%" + strconv.Itoa(percent)
		tournament.Award = total - (((ticket * quantity) * percent) / 100)
		tournament.Income = total - tournament.Award
		tournament.Sum = total
	}

	tournament.Approved = true
	tournament.State = constants.RegistrationState
	tournament.Members = 0
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

//Find all tournaments by paging
func FindAllTournaments(c *gin.Context) {

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 0, 64)
	gameKey := c.Param("game")

	gameId := driver.NewDocumentID(constants.Games, gameKey)

	//Find tournaments for game in graph by paging
	if tournamentsDoc, found := database.FindItemsInGraph(constants.GamesGraph, gameId, constants.GameToTournament, page, constants.TournamentsCount); !found {
		c.JSON(http.StatusOK, &response.Data{Data: &response.ValidationError{Error: "رقابتی برای این بازی وجود ندارد"}})
		return
	} else {

		//Get total tournaments count
		count, err := database.TotalCount(constants.GamesGraph, gameId, constants.GameToTournament)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &response.Data{Data: &response.ServerError{Message: "مشکل در دریافت تعداد رقابت ها"}})
			return
		}

		pages := utility.Pages(count, constants.TournamentsCount)

		var tournaments []map[string]interface{}

		//Remove some data for json
		for i, v := range tournamentsDoc {

			delete(tournamentsDoc[i], "_id")
			delete(tournamentsDoc[i], "_rev")
			delete(tournamentsDoc[i], "approved")
			tournamentsDoc[i]["key"] = v["_key"]
			delete(tournamentsDoc[i], "_key")

			utc, _ := time.LoadLocation("UTC")
			t := time.Unix(int64(v["created_at"].(float64)), 0)
			t = t.In(utc)
			tournamentsDoc[i]["date"] = jalali.Strftime("%A, %e %b %Y %H:%M", t)

			tournaments = append(tournaments, tournamentsDoc[i])
		}

		if tournaments == nil {
			c.JSON(http.StatusOK, &response.Data{
				Data: &response.EmptyDocument{Message: "هیچ رقابتی برای این بازی وجود ندارد"},
			})
			return
		}

		c.JSON(http.StatusOK, &response.Data{
			Data: &response.FindAllDocuments{
				Documents:   tournaments,
				TotalPages:  pages,
				CurrentPage: page,
			},
			State: true,
		})
	}

}
