package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq" // importing with name _ is special import to tell go not to remove this deps
	api "github.com/tgfukuda/be-master/api"
	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/gapi"
	"github.com/tgfukuda/be-master/pb"
	"github.com/tgfukuda/be-master/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	runGRPCServer(config, store)
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannnot create server:", err)
	}

	grpcSever := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcSever, server)
	reflection.Register(grpcSever) // add usage to server

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcSever.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannnot create server:", err)
	}

	err = server.StartServer(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start serve", err)
	}
}
