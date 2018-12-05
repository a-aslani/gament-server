package router

import (
	"Server/app/controller/v1"
	"Server/app/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {

	r := gin.Default()

	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Register", "Authorization"}
	//config.AllowAllOrigins = true
	config.AllowOrigins = []string{"http://localhost:8080"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}
	r.Use(cors.New(config))

	r.Static("/static", "./static")

	//Group V1
	v1 := r.Group("/api/v1")
	{
		//Create new phone
		v1.POST("/phone", controllerV1.GetPhoneNumber)

		//Check active code
		v1.GET("/code", controllerV1.CheckCode)
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
		//Find All Games
		v1.GET("/games", controllerV1.FindAllGames)
		//Create game
		v1.POST("/games", middleware.AdminAndOwnerAccessToken(), controllerV1.CreateGame)
		//Update game
		v1.PUT("/games/:key", middleware.AdminAndOwnerAccessToken(), controllerV1.UpdateGame)
		//Remove game
		v1.DELETE("/games/:key", middleware.AdminAndOwnerAccessToken(), controllerV1.DestroyGame)

		//Create rule
		v1.POST("/rules", controllerV1.CreateRule)

		//Create tournament
		v1.POST("/tournaments", controllerV1.CreateTournament)
	}

	return r
}
