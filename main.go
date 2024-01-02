package main

import (
	"database/sql"
	"github.com/budiharyonoo/simple_bank/api"
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln("Error connection to DB:", err)
		return
	}

	// Init sqlc
	store := db.NewStore(conn)

	// Init server
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalln("Error start the server:", err)
		return
	}
}
