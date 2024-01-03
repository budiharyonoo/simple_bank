package api

import (
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

// Server serves HTTP requests for our service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer is a contructor of Server and setup router
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Register custom validator to GIN
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatalln("Error register custom validator", err)
			return nil
		}
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is up an running!"})
	})

	// Accounts Router
	router.GET("/v1/accounts/", server.getListAccounts)
	router.POST("/v1/accounts", server.createAccount)
	router.GET("/v1/accounts/:id", server.getAccountById)

	// Transfers Router
	router.POST("/v1/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start runs HTTP Server of specific port
// address is the server address
func (server Server) Start(address string) error {
	return server.router.Run(address)
}

// paginationResponse standarized response of pagination data
// gin.H is a shortcut for map[string]any
func paginationResponse(totalRows int64, page int32, limit int32, data any) gin.H {
	return gin.H{
		"metadata": gin.H{
			"total": totalRows,
			"page":  page,
			"limit": limit,
		},
		"data": data,
	}
}

// errorResponse serialize error before send to response client
// gin.H is a shortcut for map[string]any
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
