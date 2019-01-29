package router

import (
	"Server/app/controller/v1"
	"Server/app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(*middleware.CoresConfig()))

	r.Static("/static", "./static")

	//Group V1
	v1 := r.Group("/api/v1")
	{
		//Create new phone
		v1.POST("/phone", controllerV1.GetPhoneNumber)

		//Check active code
		v1.POST("/code/:phone", controllerV1.CheckCode)
		//Create new code
		v1.POST("/code", controllerV1.RenewCode)

		//Check username
		v1.POST("/users/username", controllerV1.CheckUsername)
		//Create new user
		v1.POST("/users", middleware.RegisterAccessToken(), controllerV1.CreateUser)
		//Get user information
		v1.GET("/users/info", middleware.UserAccessToken(), controllerV1.GetUserInfo)

		//Find game
		v1.GET("/games/:key", controllerV1.FindGame)
		//Find all games
		v1.GET("/games", controllerV1.FindAllGames)
		//Create game
		v1.POST("/games", middleware.AdminAndOwnerAccessToken(), controllerV1.CreateGame)
		//Update game
		v1.PUT("/games/:key", middleware.AdminAndOwnerAccessToken(), controllerV1.UpdateGame)
		//Remove game
		v1.DELETE("/games/:key", middleware.AdminAndOwnerAccessToken(), controllerV1.DestroyGame)

		//Create rule
		v1.POST("/rules", middleware.AdminAndOwnerAccessToken(), controllerV1.CreateRule)

		//Create tournament
		v1.POST("/tournaments", middleware.AdminAndOwnerAccessToken(), controllerV1.CreateTournament)
		//Find all tournaments
		v1.GET("/tournaments/:game", controllerV1.FindAllTournaments)
	}

	return r
}
