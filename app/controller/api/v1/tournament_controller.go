package controllerApiV1

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
	count := tournament.Count

	if ticket == 0 {
		tournament.Sum = 0
		tournament.Percentage = "%0"
		tournament.Income = 0
		tournament.Award = 0
	} else {
		percent, _ := strconv.Atoi(c.PostForm("percentage"))
		total := ticket * count
		tournament.Percentage = "%" + strconv.Itoa(percent)
		tournament.Award = total - (((ticket * count) * percent) / 100)
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

		//Remove some data for json
		tournaments := refactorTournamentResponse(tournamentsDoc)

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

func refactorTournamentResponse(docs []map[string]interface{}) []map[string]interface{} {

	var newDocs []map[string]interface{}

	for i, v := range docs {

		delete(docs[i], "_id")
		delete(docs[i], "_rev")
		docs[i]["key"] = v["_key"]
		delete(docs[i], "_key")

		if _, ok := docs[i]["approved"]; ok {
			delete(docs[i], "approved")
		}

		utc, _ := time.LoadLocation("UTC")

		var t time.Time

		if docs[i]["state"] == 1 {
			t = time.Unix(int64(v["created_at"].(float64)), 0)
		} else {
			t = time.Unix(int64(v["updated_at"].(float64)), 0)
		}

		t = t.In(utc)

		docs[i]["date"] = jalali.Strftime("%A, %e %b %Y %H:%M", t)

		newDocs = append(newDocs, docs[i])
	}

	return newDocs
}
