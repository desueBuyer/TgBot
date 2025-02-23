package main

import (
	utils "tgbot/internal/utils"
	api "tgbot/pkg/api"
)

func main() {
	config := utils.InitConfig()
	api.RunBot(config)
}
