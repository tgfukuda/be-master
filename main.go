package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // importing with name _ is special import to tell go not to remove this deps
	api "github.com/tgfukuda/be-master/api"
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/util"
)

func main() {
	config, err := util.LoadConfig(".") // read
	if err != nil {
		log.Fatal("can't load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannnot create server:", err)
	}

	err = server.StartServer(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start serve", err)
	}
}
