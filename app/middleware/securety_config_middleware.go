package middleware

import "github.com/gin-contrib/cors"

func CoresConfig() *cors.Config {

	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Register", "Authorization"}
	config.AllowAllOrigins = true
	//config.AllowOrigins = []string{"http://localhost:8080", "http://192.168.1.6:1535"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}
	return &config
}