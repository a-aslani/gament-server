package utility

import (
	"fmt"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
}
