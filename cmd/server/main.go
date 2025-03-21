package main

import (
	"fmt"

	"github.com/brunohubner/golang-api/configs"
)

func main() {
	// path do env file uma pasta atrás
	config, _ := configs.LoadConfig("../")

	fmt.Println(config.DbDriver)
	fmt.Println(config.DbHost)
	fmt.Println(config.DbPort)
}
