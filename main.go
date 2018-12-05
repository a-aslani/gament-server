package main

import (
	"Server/app/constants"
	"Server/app/router"
	"Server/app/utility"
)

func main() {

	utility.InitKeys()

	r := router.New()

	r.Run(constants.Port)
}
