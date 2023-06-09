package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq" // importing with name _ is special import to tell go not to remove this deps
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	api "github.com/tgfukuda/be-master/api"
	db "github.com/tgfukuda/be-master/db/sqlc"
	_ "github.com/tgfukuda/be-master/docs/statik"
	"github.com/tgfukuda/be-master/gapi"
	"github.com/tgfukuda/be-master/pb"
	"github.com/tgfukuda/be-master/util"
	"github.com/tgfukuda/be-master/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".") // read
	if err != nil {
		log.Fatal().Err(err).Msg("can't load config")
	}

	if config.Environment == "develop" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // pretty format
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to db")
	}

	// run db migration
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServerAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)

	go runGatewayServer(config, store, taskDistributor)

	runGRPCServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migration instance")
	}

	if err := migration.Up(); err == migrate.ErrNoChange {
		log.Info().Msgf("no change detected. migration skipped")
	} else if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrate up:")
	} else if err == nil {
		log.Info().Msgf("db migrated successfully")
	}
}

func runGRPCServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannnot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcSever := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcSever, server)
	reflection.Register(grpcSever) // add usage to server

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start grpc server at %s", listener.Addr().String())
	err = grpcSever.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start grpc server")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannnot create server")
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// host swagger
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}
	// fs := http.FileServer(http.Dir("./docs/swagger")) // for directly uses js.
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start http gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start http gateway server")
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start task processor")
	}

	log.Info().Msg("start task processor")
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannnot create server")
	}

	err = server.StartServer(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot start serve")
	}
}
