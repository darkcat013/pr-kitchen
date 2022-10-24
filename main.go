package main

import (
	"math/rand"
	"time"

	"github.com/darkcat013/pr-kitchen/config"
	"github.com/darkcat013/pr-kitchen/domain"
	"github.com/darkcat013/pr-kitchen/utils"
)

func main() {

	utils.InitializeLogger()
	rand.Seed(time.Now().UnixNano())

	domain.InitializeMenu(config.MENU_PATH)
	domain.InitializeApparatuses(config.APPARATUSES_PATH)
	domain.InitializeCooks(config.COOKS_PATH)
	go domain.RunOrdersHandling()

	StartServer()
}
