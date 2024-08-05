package main

import (
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

	redisStorage, err := storage.NewRedisStorage(config.DbAddr, config.DbPass, config.Db)
	if err != nil {
		log.Fatal("Error connecting db: ", err)
	}

	server := api.NewServer(config.ListenAddr, redisStorage, config)

	log.Fatal(server.Start())
}