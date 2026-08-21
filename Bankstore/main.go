package main

import (
	"Bankstore/api"
	db "Bankstore/db/sqlc"
	"Bankstore/utils"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load configuration", err)
	}

	pool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Could not connect to db", err)
	}
	defer pool.Close()

	store := db.NewStore(pool)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Could not start server", err)
	}

}
