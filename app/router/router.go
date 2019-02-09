package router

import (
	"Server/app/controller/api/v1"
	"Server/app/controller/socket/v1"
	"Server/app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func New() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(*middleware.CoresConfig()))

	r.Static("/static", "./static")

	r.LoadHTMLFiles("index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	//Socket Group V1
	wsV1 := r.Group("/ws/v1")
	{
		//Socket connection
		wsV1.GET("/", controllerSocketV1.HandleConnection)
	}

	//API Group V1
	v1 := r.Group("/api/v1")
	{
		//Create new phone
		v1.POST("/phone", controllerApiV1.GetPhoneNumber)

		//Check active code
		v1.POST("/code/:phone", controllerApiV1.CheckCode)
		//Create new code
		v1.POST("/code", controllerApiV1.RenewCode)

		//Check username
		v1.POST("/users/username", controllerApiV1.CheckUsername)
		//Create new user
		v1.POST("/users", middleware.RegisterAccessToken(), controllerApiV1.CreateUser)
		//Get user information
		v1.GET("/users/info", middleware.UserAccessToken(), controllerApiV1.GetUserInfo)
		//Find all users
		v1.GET("/users", controllerApiV1.FindAllUsers)

		//Find game
		v1.GET("/games/:key", controllerApiV1.FindGame)
		//Find all games
		v1.GET("/games", controllerApiV1.FindAllGames)
		//Create game
		v1.POST("/games", middleware.AdminAndOwnerAccessToken(), controllerApiV1.CreateGame)
		//Update game
		v1.PUT("/games/:key", middleware.AdminAndOwnerAccessToken(), controllerApiV1.UpdateGame)
		//Remove game
		v1.DELETE("/games/:key", middleware.AdminAndOwnerAccessToken(), controllerApiV1.DestroyGame)

		//Create rule
		v1.POST("/rules", middleware.AdminAndOwnerAccessToken(), controllerApiV1.CreateRule)
		//Find rule
		v1.GET("/rules/:key", controllerApiV1.FindRule)

		//Create tournament
		v1.POST("/tournaments", middleware.AdminAndOwnerAccessToken(), controllerApiV1.CreateTournament)
		//Find all tournaments
		v1.GET("/tournaments/:game", controllerApiV1.FindAllTournaments)

		//Create team
		v1.POST("/teams", middleware.UserAccessToken(), controllerApiV1.CreateTeam)
	}

	return r
}
