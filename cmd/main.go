package main

import (
	"context"
	"os"

	"github.com/xavesen/search-admin/internal/api"
	"github.com/xavesen/search-admin/internal/config"
	"github.com/xavesen/search-admin/internal/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		os.Exit(1)
	}

	ctx := context.TODO()
	mongoStorage, err := storage.NewMongoStorage(ctx, config.DbAddr, config.Db, config.DbUser, config.DbPass)
	if err != nil {
		os.Exit(1)
	}

	server := api.NewServer(config.ListenAddr, mongoStorage, config)

	log.Fatal(server.Start())
}