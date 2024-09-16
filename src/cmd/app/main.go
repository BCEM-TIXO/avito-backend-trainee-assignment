package main

import (
	"context"
	"log"
	"tender/internal/config"
	"tender/internal/migration"
	"tender/internal/service"
	postgresql "tender/pkg/client"
)

func main() {
	cfg := config.GetConfig()
	pSQLClinet, _ := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	err := migration.RunMigrations(pSQLClinet)
	if err != nil {
		log.Fatal("No DB connection ", err.Error())
	}
	service := service.NewService(pSQLClinet, cfg.ServerAddress)
	service.Run()
}
