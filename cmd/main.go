package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lailiseptiandi/be-kominfo-test/internal/routers"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/config"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/database"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dbs, err := database.InitDB(*cfg)

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	defer database.DisconnectDB(dbs)

	r := gin.Default()
	routers.InitRoutes(dbs, r)
	serverAdd := fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)
	r.Run(serverAdd)

	fmt.Println("Success starting serve")
}
