package main

import (
	"database/sql"
	"github.com/budiharyonoo/simple_bank/api"
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/budiharyonoo/simple_bank/utils"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalln("Error read config file:", err)
		return
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("Error connection to DB:", err)
		return
	}

	// Init sqlc
	store := db.NewStore(conn)

	// Init server
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("Error start the server:", err)
		return
	}
}
