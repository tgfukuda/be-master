package gapi

import (
	"fmt"

	db "github.com/tgfukuda/be-master/db/sqlc"
	"github.com/tgfukuda/be-master/pb"
	"github.com/tgfukuda/be-master/token"
	"github.com/tgfukuda/be-master/util"
	"github.com/tgfukuda/be-master/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer // for forward compatibility
	config                           util.Config
	store                            db.Store
	tokenMaker                       token.Maker
	taskDistributor                  worker.TaskDistributor
}

// new Http Server and setup routes
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
