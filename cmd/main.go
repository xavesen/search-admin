package main

import (
	"context"
	"log"

	"github.com/xavesen/search-admin/internal/api"
	"github.com/xavesen/search-admin/internal/config"
	"github.com/xavesen/search-admin/internal/storage"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error getting config from evironment: ", err)
	}

	ctx := context.TODO()
	mongoStorage, err := storage.NewMongoStorage(ctx, config.DbAddr, config.Db, config.DbUser, config.DbPass)
	if err != nil {
		log.Fatal("Error connecting db: ", err)
	}

	server := api.NewServer(config.ListenAddr, mongoStorage, config)

	log.Fatal(server.Start())
}