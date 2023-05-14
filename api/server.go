package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/tgfukuda/be-master/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// new Http Server and setup routes
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccount)

	server.router = router

	return server
}

// Start runs HTTP server on a specific address. - the reason why implement it in such a way is router is private and we have a gracefully shutdown logic in this.
func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
