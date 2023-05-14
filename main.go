package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // importing with name _ is special import to tell go not to remove this deps
	api "github.com/tgfukuda/be-master/api"
	db "github.com/tgfukuda/be-master/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	var err error
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer((store))

	err = server.StartServer(serverAddress)

	if err != nil {
		log.Fatal("cannot start serve", err)
	}
}
