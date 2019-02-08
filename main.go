package main

import (
	"Server/app/constants"
	"Server/app/router"
)

func main() {



	//utility.InitKeys()

	r := router.New()

	r.Run(constants.Port)
}
