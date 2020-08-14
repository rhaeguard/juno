package main

import (
	_ "github.com/lib/pq"
	"github.com/mensurowary/juno/config"
	database "github.com/mensurowary/juno/db"
	"github.com/mensurowary/juno/router"
)

func main() {
	db := database.Initialize()
	defer db.Close()

	engine := router.Initialize(db)
	_ = engine.Run(":" + config.Config.Port)
}
